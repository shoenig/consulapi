package consulapi

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type SessionID string
type SessionTerminationBehavior string

const (
	SessionRelease SessionTerminationBehavior = "release"
	SessionDelete  SessionTerminationBehavior = "delete"

	SessionMinimumTTL = 10 * time.Second
	SessionMaximumTTL = 86400 * time.Second

	SessionMinimumLockDelay = 0 * time.Second
	SessionMaximumLockDelay = 60 * time.Second
)

type SessionConfig struct {
	// The DC in which the node holding the session.
	DC string `json:"-"` // not part of the official API

	// The node with which the session is associated. Typically, this should be
	// set to the node name of the local consul agent. That information can be
	// retrieved using the Self endpoint.
	Node string `json:"Node"`

	// Name is a human-readable identifier for this session.
	Name string `json:"Name"`

	// LockDelay determines the minimum amount of time that must pass
	// after a lock expiration before a new lock acquisition may take place.
	// The use of such a delay is to allow the potentially still-alive leader
	// to detect its own lock invalidation and stop processing requests that
	// may lead to an inconsistent state. This mechanism is not "bullet-proof",
	// but it eliminates the need for clients to include their own timeout code.
	// A zero-value disables this feature.
	LockDelay time.Duration `json:"LockDelay"`

	// TTL represents the amount of time that may pass before the session
	// automatically becomes invalidated. Typically a lower value for TTL
	// is better, as a node that unexpectedly goes "off the grid" will
	// continue to own the lock until the TTL expires.
	TTL time.Duration `json:"TTL"`

	// Behavior controls what action to take when a session is invalidated.
	// - SessionRelease: causes any held locks to be released.
	// - SessionDelete: causes any held locks to be deleted.
	Behavior SessionTerminationBehavior `json:"Behavior"`
}

type sessionConfigFormat2 struct {
	Node      string `json:"Node"`
	Name      string `json:"Name"`
	LockDelay string `json:"LockDelay"`
	TTL       string `json:"TTL"`
	Behavior  string `json:"Behavior"`
}

type sessionConfigFormat3 struct {
	ID        string  `json:"ID"`
	Node      string  `json:"Node"`
	Name      string  `json:"Name"`
	LockDelay float64 `json:"LockDelay"`
	TTL       string  `json:"TTL"`
	Behavior  string  `json:"Behavior"`
}

// SessionQuery is used to define values for each of the optional parameters
// on the session endpoint.
type SessionQuery struct {
	// ID of the session for which the query is regarding.
	ID SessionID

	// DC indicates which dc to query.
	//
	// If blank, this will default to the dc that the queried agent is in.
	DC string
}

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Session -s _mock.go

type Session interface {
	// CreateSession will initialize a new session.
	//
	// https://www.consul.io/api/session.html#create-session
	CreateSession(Ctx, SessionConfig) (SessionID, error)

	// DeleteSession will delete session of id.
	//
	// https://www.consul.io/api/session.html#delete-session
	DeleteSession(Ctx, SessionQuery) error

	// ReadSession will return the session information for id.
	//
	// https://www.consul.io/api/session.html#read-session
	ReadSession(Ctx, SessionQuery) (SessionConfig, error)

	// ListSessions will list every session on node.
	//
	// https://www.consul.io/api/session.html#list-sessions-for-node
	ListSessions(ctx Ctx, dc, node string) (map[SessionID]SessionConfig, error)

	// RenewSession will renew session of id.
	//
	// https://www.consul.io/api/session.html#renew-session
	RenewSession(Ctx, SessionQuery) (time.Duration, error)
}

var _ Session = (*client)(nil)

func internalizeSessionConfig(config SessionConfig) (sessionConfigFormat2, error) {
	var isc sessionConfigFormat2

	if config.Node == "" {
		return isc, errors.New("session node required")
	}
	isc.Node = config.Node

	if config.Name == "" {
		return isc, errors.New("session name required")
	}
	isc.Name = config.Name

	if (config.LockDelay < SessionMinimumLockDelay) || (config.LockDelay > SessionMaximumLockDelay) {
		return isc, errors.New("session lock delay must be at least 0 but at most 60 seconds")
	}
	isc.LockDelay = config.LockDelay.String()

	if (config.TTL < SessionMinimumTTL) || (config.TTL > SessionMaximumTTL) {
		return isc, errors.New("session ttl must be more than 10 but less than 86400 seconds")
	}
	isc.TTL = config.TTL.String()

	if (config.Behavior != SessionRelease) && (config.Behavior != SessionDelete) {
		return isc, errors.New("session behavior must be 'release' or 'delete'")
	}
	isc.Behavior = string(config.Behavior)

	return isc, nil
}

func (c *client) CreateSession(ctx Ctx, config SessionConfig) (SessionID, error) {
	dc := config.DC

	isc, err := internalizeSessionConfig(config)
	if err != nil {
		return "", err
	}

	path := fixup("/v1/session", "create", param("dc", dc))

	var response = struct {
		ID SessionID `json:"ID"`
	}{}

	body, err := json.Marshal(isc)
	if err != nil {
		return "", err
	}

	if err := c.put(ctx, path, string(body), &response); err != nil {
		return "", errors.Wrap(err, "failed to create session")
	}

	return response.ID, nil
}

func (c *client) ReadSession(ctx Ctx, query SessionQuery) (SessionConfig, error) {
	id := query.ID
	dc := query.DC

	path := fixup("/v1/session/info", string(id), param("dc", dc))

	var response []sessionConfigFormat3
	if err := c.get(ctx, path, &response); err != nil {
		return SessionConfig{}, errors.Wrap(err, "failed to read session")
	}

	return sessionFromFormat3(response, dc)
}

func (c *client) RenewSession(ctx Ctx, query SessionQuery) (time.Duration, error) {
	id := query.ID
	dc := query.DC

	path := fixup("/v1/session/renew", string(id), param("dc", dc))

	var response []sessionConfigFormat3
	if err := c.put(ctx, path, "", &response); err != nil {
		return 0, errors.Wrap(err, "failed to renew session")
	}

	session, err := sessionFromFormat3(response, dc)
	return session.TTL, err
}

func sessionFromFormat3(response []sessionConfigFormat3, dc string) (SessionConfig, error) {
	if len(response) < 1 {
		return SessionConfig{}, errors.New("read session returned no sessions")
	}

	if len(response) > 1 {
		return SessionConfig{}, errors.New("read session returned more than one session")
	}

	session := response[0]

	ttl, err := time.ParseDuration(session.TTL)
	if err != nil {
		return SessionConfig{}, err
	}

	return SessionConfig{
		DC:        dc,
		Node:      session.Node,
		Name:      session.Name,
		LockDelay: time.Duration(session.LockDelay),
		TTL:       ttl,
		Behavior:  SessionTerminationBehavior(session.Behavior),
	}, nil
}

func (c *client) ListSessions(ctx Ctx, dc, node string) (map[SessionID]SessionConfig, error) {
	path := fixup("/v1/session/node/", node, param("dc", dc))

	var response []sessionConfigFormat3
	if err := c.get(ctx, path, &response); err != nil {
		return nil, errors.Wrap(err, "failed to list sessions")
	}

	configs := make(map[SessionID]SessionConfig, len(response))

	for _, session := range response {
		ttl, err := time.ParseDuration(session.TTL)
		if err != nil {
			return nil, err
		}

		configs[SessionID(session.ID)] = SessionConfig{
			Node:      session.Node,
			Name:      session.Name,
			LockDelay: time.Duration(session.LockDelay),
			TTL:       ttl,
			Behavior:  SessionTerminationBehavior(session.Behavior),
		}
	}

	return configs, nil
}

func (c *client) DeleteSession(ctx Ctx, query SessionQuery) error {
	id := query.ID
	dc := query.DC

	path := fixup("/v1/session/destroy", string(id), param("dc", dc))

	if err := c.put(ctx, path, "", nil); err != nil {
		return errors.Wrap(err, "failed to destroy session")
	}

	return nil
}
