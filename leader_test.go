package consulapi

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type leadershipHandler struct {
	t *testing.T
}

func (lh *leadershipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/v1/agent/self":
		lh.rxSelf(w)
	case "/v1/kv/my/key":
		lh.rxMyKey(w)
	case "/v1/session/create":
		lh.rxCreate(w)
	default:
		lh.t.Fatal("unexpected path:", r.URL.Path)
	}
}

func (lh *leadershipHandler) rxSelf(w http.ResponseWriter) {
	response := load(lh.t, "test_leadership_self.json")
	_, _ = io.WriteString(w, response)
}

func (lh *leadershipHandler) rxMyKey(w http.ResponseWriter) {
	response := load(lh.t, "test_leadership_myKey.json")
	_, _ = io.WriteString(w, response)
}

func (lh *leadershipHandler) rxCreate(w http.ResponseWriter) {
	response := load(lh.t, "test_leadership_create.json")
	_, _ = io.WriteString(w, response)
}

func Test_Leadership_Participate(t *testing.T) {

	ctx, ts, client := testClient(&leadershipHandler{
		t: t,
	})
	defer ts.Close()

	var alf AsLeaderFunc = func(Ctx) error {
		t.Log("this is as leader func")
		return nil
	}

	session, err := client.Participate(ctx, LeadershipConfig{
		Key:         "/my/key",
		ContactInfo: "node1",
		Description: "the session key",
		LockDelay:   5 * time.Second,
		TTL:         30 * time.Second,
	}, alf)

	require.NoError(t, err)
	t.Log("session:", session)

	current, err := session.Current(ctx)
	require.NoError(t, err)
	t.Log("current:", current)

	// i only feel a little bad about this
	for i := 0; wait(t, i, session); i++ {
		time.Sleep(10 * time.Millisecond)
	}
}

func wait(t *testing.T, i int, session LeaderSession) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	id := session.SessionID(ctx)
	t.Logf("waiting for session [%d]: %s", i, id)
	return id == ""
}
