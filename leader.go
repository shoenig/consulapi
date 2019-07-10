package consulapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type LeadershipConfig struct {
	// Key is the path used to store session information in
	// the consul KV store. This must be set for the client to be used for
	// leader election. Typically this value will look something like
	// "service/<service name>/leader".
	Key string

	// ContactInfo is an opaque string that identifies the leader elected node.
	// This is what clients should use to know how to connect to the service
	// that is currently the leader. How this string is created, parsed, and
	// interpreted is an implementation detail left to the clients. Typically,
	// this string will in the form of a URI.
	ContactInfo string

	// Description is an arbitrary human-readable name for this leader election.
	// If not set, Description defaults to "default-leader-session".
	Description string

	// See SessionConfig.LockDelay.
	LockDelay time.Duration

	// See SessionConfig.TTL.
	TTL time.Duration
}

func (lc LeadershipConfig) name() string {
	if lc.Description == "" {
		return "default-leader-session"
	}
	return lc.Description
}

// AsLeaderFunc is executed when the client is able to acquire
// the underlying consul leader lock. When a value is sent on context.Done,
// this client is no longer the elected leader and must cease any operations
// that require leadership. A reasonable implementation will return from the
// function as soon as possible context cancellation.
//
// If the implementation returns from the function before the context is
// cancelled, leadership will be abdicated, the context will be cancelled,
// and no further action should be taken until elected leader again.
type AsLeaderFunc func(Ctx) error

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Candidate -s _mock.go

// A Candidate implementation is able to Participate in leadership elections.
type Candidate interface {
	Participate(Ctx, LeadershipConfig, AsLeaderFunc) (LeaderSession, error)
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i LeaderSession -s _mock.go

// A LeaderSession is used to inspect and manage the underlying consul session
// being used to participate in leadership elections.
type LeaderSession interface {
	Abdicate(Ctx) error
	Current(Ctx) (string, error)
	SessionID(Ctx) string
}

type leadershipManager struct {
	client *client
	key    string

	self        AgentInfo
	contactInfo string
	sessionTTL  time.Duration
	sessionID   atomic.Value
	isLeader    atomic.Value

	asLeader AsLeaderFunc
}

func (c *client) Participate(ctx Ctx, opts LeadershipConfig, f AsLeaderFunc) (LeaderSession, error) {
	self, err := c.Self(ctx)
	if err != nil {
		return nil, err
	}

	manager := &leadershipManager{
		client: c,
		key:    strings.TrimPrefix(opts.Key, "/"),

		self:        self,
		contactInfo: opts.ContactInfo,
		sessionTTL:  opts.TTL,
		asLeader:    f,
	}

	go manager.maintainSession(opts)
	go manager.maintainLeadership()

	return manager, nil
}

func (lm *leadershipManager) maintainSession(opts LeadershipConfig) {
	for {

		if err := lm.createSession(opts); err != nil {
			lm.client.log.Warnf("failed to create session, try again in 3s: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		// (lock delay covers the fence post)
		for range time.Tick(lm.sessionTTL) {
			ctx := context.TODO()
			if _, err := lm.client.RenewSession(ctx, SessionQuery{
				DC: "",
				ID: lm.getSessionID(),
			}); err != nil {
				lm.client.log.Warnf("failed to renew session, will need to create a new one")
				break
			}
		}
	}
}

func (lm *leadershipManager) createSession(opts LeadershipConfig) error {
	lm.client.log.Tracef("attempting to establish new session")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sessionID, err := lm.client.CreateSession(ctx, SessionConfig{
		// DC not used
		Node:      lm.self.Name,
		Name:      opts.name(),
		LockDelay: opts.LockDelay,
		TTL:       opts.TTL,
		Behavior:  SessionDelete,
	})
	if err != nil {
		return err
	}

	lm.setSessionID(sessionID)
	lm.client.log.Tracef("established new session with ID: %s", sessionID)
	return nil
}

func (lm *leadershipManager) setSessionID(id SessionID) {
	lm.sessionID.Store(id)
}

func (lm *leadershipManager) getSessionID() SessionID {
	value := lm.sessionID.Load()
	if value == nil {
		return SessionID("")
	}
	return value.(SessionID)
}

func (lm *leadershipManager) Abdicate(ctx Ctx) error {
	lm.isLeader.Store(false)

	path := fixup("/v1/kv/", lm.key, param("release", string(lm.getSessionID())))

	var response bool
	if err := lm.client.put(ctx, path, lm.value(), &response); err != nil {
		return errors.Wrap(err, "failed to abdicate leadership")
	}

	return nil
}

func (lm *leadershipManager) Current(ctx Ctx) (string, error) {
	path := fixup("/v1/kv/", lm.key)

	var response []struct {
		Value string `json:"value"`
	}

	if err := lm.client.get(ctx, path, &response); err != nil {
		return "", errors.Wrap(err, "failed to lookup leadership")
	}

	encoded := response[0].Value
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode leadership identity")
	}

	return string(decoded), nil
}

func (lm *leadershipManager) SessionID(_ Ctx) string {
	return string(lm.getSessionID())
}

func (lm *leadershipManager) value() string {
	bs, err := json.Marshal(lm.contactInfo)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

func (lm *leadershipManager) tryAcquire() (bool, error) {
	id := string(lm.getSessionID())
	if id == "" {
		return false, errors.New("cannot acquire leader lock before establishing session")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	path := fixup("/v1/kv/", lm.key, param("acquire", id))
	var response bool
	if err := lm.client.put(ctx, path, lm.value(), &response); err != nil {
		return false, errors.Wrap(err, "failed to acquire leadership")
	}

	return response, nil
}

func (lm *leadershipManager) maintainLeadership() {
	// initially we are not the leader
	lm.setIsLeader(false)

	// initial gap is very small so we can get started now
	gap := 1 * time.Millisecond

	for {
		time.Sleep(gap)

		// try to acquire leadership
		won, err := lm.tryAcquire()
		switch {
		case err != nil:
			gap = 1 * time.Second
			lm.client.log.Warnf("encountered error while acquiring leadership, try again in 1 second: %v", err)
			continue

		case !won:
			gap = lm.sessionTTL
			continue

		case won:
			ctx, cancel := context.WithCancel(context.Background())

			// in the background, try to maintain leadership until we lose it
			go func(ctx context.Context) {
				ticker := time.NewTicker(lm.sessionTTL)
				for {
					select {
					case <-ticker.C:
						won, err := lm.tryAcquire()
						if err != nil || !won {
							cancel()
							return
						}
					case <-ctx.Done():
						return
					}
				}
			}(ctx)

			// in the foreground, run the AsLeaderFunc until it returns or ctx is cancelled
			err := lm.asLeader(ctx)
			lm.client.log.Errorf("provided AsLeaderFunc returned with error: %v", err)
			lm.setIsLeader(false)
			cancel() // probably why it returned, but we still need to run it in case
		}
	}
}

func (lm *leadershipManager) IsLeader() bool {
	isLeader, ok := lm.isLeader.Load().(bool)
	if !ok {
		return false
	}
	return isLeader
}

func (lm *leadershipManager) setIsLeader(b bool) {
	lm.client.log.Tracef("setting leader status to: %v", b)
	lm.isLeader.Store(b)
}

func (lm *leadershipManager) String() string {
	return fmt.Sprintf("[leader(%s):%t]", lm.key, lm.IsLeader())
}
