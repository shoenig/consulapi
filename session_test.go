package consulapi

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Session_CreateSession(t *testing.T) {
	expPayload := `{"Node":"dc1-node1","Name":"mySession1","LockDelay":"1s","TTL":"10s","Behavior":"release"}`
	expID := SessionID("adf4238a-882b-9ddc-4a9d-5b6758e4159e")

	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_create.json"),
		hasPath:   "/v1/session/create",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{},
		hasBody:   expPayload,
	})
	defer ts.Close()

	id, err := client.CreateSession(ctx, SessionConfig{
		// DC not set
		Node:      "dc1-node1",
		Name:      "mySession1",
		LockDelay: 1 * time.Second,
		TTL:       10 * time.Second,
		Behavior:  SessionRelease,
	})
	require.NoError(t, err)
	require.Equal(t, expID, id)
}

func Test_Session_CreateSession_dc(t *testing.T) {
	expPayload := `{"Node":"dc2-node1","Name":"mySession1","LockDelay":"1s","TTL":"10s","Behavior":"release"}`
	expID := SessionID("adf4238a-882b-9ddc-4a9d-5b6758e4159e")

	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_create.json"),
		hasPath:   "/v1/session/create",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"dc": {"dc2"},
		},
		hasBody: expPayload,
	})
	defer ts.Close()

	id, err := client.CreateSession(ctx, SessionConfig{
		DC:        "dc2",
		Node:      "dc2-node1",
		Name:      "mySession1",
		LockDelay: 1 * time.Second,
		TTL:       10 * time.Second,
		Behavior:  SessionRelease,
	})
	require.NoError(t, err)
	require.Equal(t, expID, id)
}

func Test_Session_CreateSession_err(t *testing.T) {
	expPayload := `{"Node":"dc1-node1","Name":"mySession1","LockDelay":"1s","TTL":"10s","Behavior":"release"}`

	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/session/create",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{},
		hasBody:   expPayload,
	})
	defer ts.Close()

	_, err := client.CreateSession(ctx, SessionConfig{
		Node:      "dc1-node1",
		Name:      "mySession1",
		LockDelay: 1 * time.Second,
		TTL:       10 * time.Second,
		Behavior:  SessionRelease,
	})
	require.EqualError(t, err, "failed to create session: status code (500)")
}

func Test_Session_DeleteSession(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_delete.json"),
		hasPath:   "/v1/session/destroy/abc123",
		hasMethod: http.MethodPut,
	})
	defer ts.Client()

	err := client.DeleteSession(ctx, SessionQuery{
		// no DC set
		ID: "abc123",
	})
	require.NoError(t, err)
}

func Test_Session_DeleteSession_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_delete.json"),
		hasPath:   "/v1/session/destroy/abc123",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"dc": {"dc2"},
		},
	})
	defer ts.Client()

	err := client.DeleteSession(ctx, SessionQuery{
		DC: "dc2",
		ID: "abc123",
	})
	require.NoError(t, err)
}

func Test_Session_DeleteSession_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/session/destroy/abc123",
		hasMethod: http.MethodPut,
	})
	defer ts.Client()

	err := client.DeleteSession(ctx, SessionQuery{
		ID: "abc123",
	})
	require.EqualError(t, err, "failed to destroy session: status code (500)")
}

func Test_Session_ReadSession(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_info.json"),
		hasPath:   "/v1/session/info/abc123",
		hasMethod: http.MethodGet,
	})
	defer ts.Client()

	config, err := client.ReadSession(ctx, SessionQuery{
		// No DC
		ID: "abc123",
	})
	require.NoError(t, err)
	require.Equal(t, "test-session", config.Name)
}

func Test_Session_ReadSession_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_info.json"),
		hasPath:   "/v1/session/info/abc123",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc": {"dc2"},
		},
	})
	defer ts.Client()

	config, err := client.ReadSession(ctx, SessionQuery{
		DC: "dc2",
		ID: "abc123",
	})
	require.NoError(t, err)
	require.Equal(t, "test-session", config.Name)
	require.Equal(t, "dc2", config.DC)
}

func Test_Session_ReadSession_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/session/info/abc123",
		hasMethod: http.MethodGet,
	})
	defer ts.Client()

	_, err := client.ReadSession(ctx, SessionQuery{
		ID: "abc123",
	})
	require.EqualError(t, err, "failed to read session: status code (500)")
}

func Test_Session_RenewSession(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_renew.json"),
		hasPath:   "/v1/session/renew/abc123",
		hasMethod: http.MethodPut,
	})
	defer ts.Client()

	ttl, err := client.RenewSession(ctx, SessionQuery{
		// No DC
		ID: "abc123",
	})
	require.NoError(t, err)
	require.Equal(t, 15*time.Second, ttl)
}

func Test_Session_RenewSession_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_renew.json"),
		hasPath:   "/v1/session/renew/abc123",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"dc": {"dc2"},
		},
	})
	defer ts.Client()

	ttl, err := client.RenewSession(ctx, SessionQuery{
		DC: "dc2",
		ID: "abc123",
	})
	require.NoError(t, err)
	require.Equal(t, 15*time.Second, ttl)
}

func Test_Session_RenewSession_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/session/renew/abc123",
		hasMethod: http.MethodPut,
	})
	defer ts.Client()

	_, err := client.RenewSession(ctx, SessionQuery{
		ID: "abc123",
	})
	require.EqualError(t, err, "failed to renew session: status code (500)")
}

func Test_Session_ListSession(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_session_list.json"),
		hasPath:   "/v1/session/node/node1",
		hasMethod: http.MethodGet,
	})
	defer ts.Client()

	m, err := client.ListSessions(ctx, "", "node1")
	require.NoError(t, err)

	a, okA := m["2152c206-ccb9-5451-afb8-f3a9f5c2a374"]
	b, okB := m["35e63278-0310-c915-6a4a-3f8f27f3a02d"]
	c, okC := m["492c9bea-ce34-8386-44de-1f9105b8a0d5"]

	require.True(t, okA)
	require.True(t, okB)
	require.True(t, okC)

	require.Equal(t, "common/consul/lock/tst-node1", a.Name)
	require.Equal(t, "testbird-mon-leader-session", b.Name)
	require.Equal(t, "", c.Name)
}
