package consulapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests are for the /v1/agent/* endpoints
//
// https://www.consul.io/api/agent.html
//
// Not supported:
// - members with ?segment parameter
// - stream logs
// - update acl tokens

func Test_Client_v1_agent_self_ok(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_self.json"),
		hasPath:   "/v1/agent/self",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	self, err := client.Self(ctx)
	require.NoError(t, err)
	require.Equal(t, "mydc-mynode1", self.Name)
	require.Equal(t, "10.3.0.19", self.Address)
	require.Equal(t, "mydc", self.Tags["dc"])
}

func Test_Client_v1_agent_self_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/self",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	_, err := client.Self(ctx)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_members(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_members.json"),
		hasPath:   "/v1/agent/members",
		hasMethod: http.MethodGet,
		// hasQuery: map[string]string{"wan": "false"}, // nope; see impl
	})
	defer ts.Close()

	members, err := client.Members(ctx, false)
	require.NoError(t, err)
	require.Equal(t, 32, len(members))
	require.Equal(t, "mydc-systems1", members[0].Name)
}

func Test_Client_v1_agent_members_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/members",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	_, err := client.Members(ctx, false)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_members_wan(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_members-wan.json"),
		hasPath:   "/v1/agent/members",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{"wan": {"true"}},
	})

	defer ts.Close()

	members, err := client.Members(ctx, true)
	require.NoError(t, err)
	require.Equal(t, 1, len(members))
	require.Equal(t, "otherdc-systems1.otherdc", members[0].Name)
}

func Test_Client_v1_agent_members_wan_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/members",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{"wan": {"true"}},
	})
	defer ts.Close()

	_, err := client.Members(ctx, true)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_reload(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_reload.json"),
		hasPath:   "/v1/agent/reload",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	err := client.Reload(ctx)
	require.NoError(t, err)
}

func Test_Client_v1_agent_reload_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/reload",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	err := client.Reload(ctx)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_maintenance_enable(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_maintenance-enable.json"),
		hasPath:   "/v1/agent/maintenance",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"enable": {"true"},
			"reason": {"my reason"},
		},
	})
	defer ts.Close()

	err := client.MaintenanceMode(ctx, true, "my reason")
	require.NoError(t, err)
}

func Test_Client_v1_agent_maintenance_enable_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/maintenance",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"enable": {"true"},
			"reason": {"my reason"},
		},
	})
	defer ts.Close()

	err := client.MaintenanceMode(ctx, true, "my reason")
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_maintenance_disable(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_maintenance-disable.json"),
		hasPath:   "/v1/agent/maintenance",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"enable": {"false"},
			"reason": {"my reason"},
		},
	})
	defer ts.Close()

	err := client.MaintenanceMode(ctx, false, "my reason")
	require.NoError(t, err)
}

func Test_Client_v1_agent_maintenance_disable_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/maintenance",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"enable": {"false"},
			"reason": {"my reason"},
		},
	})
	defer ts.Close()

	err := client.MaintenanceMode(ctx, false, "my reason")
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_metrics(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_metrics.json"),
		hasPath:   "/v1/agent/metrics",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	metrics, err := client.Metrics(ctx)
	require.NoError(t, err)
	require.Equal(t, "2019-07-14 17:49:10 +0000 UTC", metrics.Timestamp)
	require.Equal(t, "mydc-mynode1", metrics.Counters[0].Labels["node"])
}

func Test_Client_v1_agent_metrics_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/metrics",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	_, err := client.Metrics(ctx)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_join(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_join.json"),
		hasPath:   "/v1/agent/join/10.0.0.1",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{"wan": {"false"}},
	})
	defer ts.Close()

	err := client.Join(ctx, "10.0.0.1", false)
	require.NoError(t, err)
}

func Test_Client_v1_agent_join_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/join/10.0.0.1",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{"wan": {"false"}},
	})
	defer ts.Close()

	err := client.Join(ctx, "10.0.0.1", false)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_join_wan(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_join.json"),
		hasPath:   "/v1/agent/join/10.0.0.1",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{"wan": {"true"}},
	})
	defer ts.Close()

	err := client.Join(ctx, "10.0.0.1", true)
	require.NoError(t, err)
}

func Test_Client_v1_agent_join_wan_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/join/10.0.0.1",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{"wan": {"true"}},
	})
	defer ts.Close()

	err := client.Join(ctx, "10.0.0.1", true)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_leave(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_leave.json"),
		hasPath:   "/v1/agent/leave",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	err := client.Leave(ctx)
	require.NoError(t, err)
}

func Test_Client_v1_agent_leave_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/leave",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	err := client.Leave(ctx)
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_forceLeave(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_agent_forceLeave.json"),
		hasPath:   "/v1/agent/force-leave/badNode1",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	err := client.ForceLeave(ctx, "badNode1")
	require.NoError(t, err)
}

func Test_Client_v1_agent_forceLeave_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/agent/force-leave/badNode1",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	err := client.ForceLeave(ctx, "badNode1")
	require.EqualError(t, err, "status code (500)")
}

func Test_Client_v1_agent_token_unknown_kind(t *testing.T) {
	ctx, ts, client := testClient(&responder{t: t})
	defer ts.Close()

	err := client.SetACLToken(ctx, "badKind", "abc123")
	require.EqualError(t, err, `unrecognized kind of token "badKind"`)
}

func Test_Client_v1_agent_token_default(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "",
		hasPath:   "/v1/agent/token/default",
		hasMethod: http.MethodPut,
		hasBody:   `{"Token":"abc123"}`,
	})
	defer ts.Close()

	err := client.SetACLToken(ctx, "default", "abc123")
	require.NoError(t, err)
}

func Test_Client_v1_agent_token_default_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "",
		hasPath:   "/v1/agent/token/default",
		hasMethod: http.MethodPut,
		hasBody:   `{"Token":"abc123"}`,
	})
	defer ts.Close()

	err := client.SetACLToken(ctx, "default", "abc123")
	require.EqualError(t, err, "status code (500)")
}
