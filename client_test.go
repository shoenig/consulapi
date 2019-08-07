package consulapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fixup(t *testing.T) {
	tests := []struct {
		prefix string
		path   string
		params [][2]string
		exp    string
	}{
		{prefix: "/v1/kv", path: "/devops/k1", exp: "/v1/kv/devops/k1"},
		{prefix: "/v1/kv", path: "/devops/k2", params: [][2]string{{"dc", "aus"}}, exp: "/v1/kv/devops/k2?dc=aus"},
		{prefix: "/v1/kv", path: "/devops/k3", params: [][2]string{{"dc", "aus"}, {"list", "true"}}, exp: "/v1/kv/devops/k3?dc=aus&list=true"},
	}

	for _, test := range tests {
		fixed := fixup(test.prefix, test.path, test.params...)
		t.Log("fixed", fixed)
		require.Equal(t, test.exp, fixed)
	}
}

func Test_param(t *testing.T) {
	k, v := "k", "v"
	pair := param(k, v)
	require.Equal(t, "k", pair[0])
	require.Equal(t, "v", pair[1])
}

func Test_RequestError_StatusCode(t *testing.T) {
	re := RequestError{
		statusCode: http.StatusTeapot,
	}

	code := re.StatusCode()
	require.Equal(t, http.StatusTeapot, code)
}

func Test_Client_New_defaults(t *testing.T) {
	c := New(ClientOptions{
		// empty, use defaults
	}).(*client)

	require.Equal(t, "http://localhost:8500", c.address)
	require.Equal(t, "", c.token)
	require.NotNil(t, c.httpClient)
	require.NotNil(t, c.log)
}

type myFoo struct {
	Foo string `json:"foo"`
}

func Test_Client_get(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "test_myfoo.json"),
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).get(ctx, "/test/arbitrary", &value)
	require.NoError(t, err)
	require.Equal(t, "bar", value.Foo)
}

func Test_Client_get_bad_request(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "test_myfoo.json"),
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).get(ctx, "not_a_path", &value)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid URL port")
}

func Test_Client_get_bad_reply(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "",
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodGet,
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).get(ctx, "/test/arbitrary", &value)
	require.EqualError(t, err, "EOF")
}

func Test_Client_get_bad_response(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusTeapot,
		body:      load(t, "test_myfoo.json"),
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).get(ctx, "/test/arbitrary", &value)
	require.EqualError(t, err, "status code (418)")
}

const (
	egBody = `{"a":1}`
)

func Test_Client_put(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "test_myfoo.json"),
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodPut,
		hasBody:   egBody,
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).put(ctx, "/test/arbitrary", egBody, &value)
	require.NoError(t, err)
	require.Equal(t, "bar", value.Foo)
}

func Test_Client_put_bad_request(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "",
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodPut,
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).put(ctx, "not_a_path", egBody, &value)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid URL port")
}

func Test_Client_put_bad_response(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusTeapot,
		body:      load(t, "test_myfoo.json"),
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodPut,
		hasBody:   egBody,
	})
	defer ts.Close()

	var value myFoo
	err := c.(*client).put(ctx, "/test/arbitrary", egBody, &value)
	require.EqualError(t, err, "status code (418)")
}

func Test_Client_delete(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "",
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodDelete,
	})
	defer ts.Close()

	err := c.(*client).delete(ctx, "/test/arbitrary")
	require.NoError(t, err)
}

func Test_Client_delete_bad_request(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "",
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodDelete,
	})
	defer ts.Close()

	err := c.(*client).delete(ctx, "not_a_path")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid URL port")
}

func Test_Client_delete_bad_response(t *testing.T) {
	ctx, ts, c := testClient(&responder{
		t:         t,
		code:      http.StatusTeapot,
		body:      "malfunction",
		hasPath:   "/test/arbitrary",
		hasMethod: http.MethodDelete,
	})
	defer ts.Close()

	err := c.(*client).delete(ctx, "/test/arbitrary")
	require.EqualError(t, err, "status code (418)")
}
