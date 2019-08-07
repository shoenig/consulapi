package consulapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Get Put Delete Keys Recurse

/*
Hashicorp has a nice demo setup to play around with, e.g.

$ curl -sL "https://demo.consul.io/v1/kv/config?keys=true" | jq .
[
  "config/",
  "config/baz/",
  "config/baz/mysql-host",
  "config/baz/xxx",
  "config/foo",
  "config/name",
  "config/remi",
  "config/test"
]
*/

func Test_KV_Get(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_kv_config_baz_bar.json"),
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	v, err := client.Get(ctx, "config/baz/bar", Query{})
	require.NoError(t, err)
	require.Equal(t, "myValue", v)
}

func Test_KV_Get_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_kv_config_baz_bar.json"),
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"dc": {"dc1"},
		},
	})
	defer ts.Close()

	v, err := client.Get(ctx, "config/baz/bar", Query{DC: "dc1"})
	require.NoError(t, err)
	require.Equal(t, "myValue", v)
}

func Test_KV_Get_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      load(t, "v1_kv_config_baz_bar.json"),
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Get(ctx, "config/baz/bar", Query{})
	require.EqualError(t, err, "status code (500)")
}

func Test_KV_Get_non_existent(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusNotFound,
		body:      "",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodGet,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	_, err := client.Get(ctx, "config/baz/bar", Query{})
	require.EqualError(t, err, `key "/v1/kv/config/baz/bar" does not exist`)
}

func Test_KV_Put(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "true",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{},
		hasBody:   "someValue",
	})
	defer ts.Close()

	err := client.Put(ctx, "config/baz/bar", "someValue", Query{})
	require.NoError(t, err)
}

func Test_KV_Put_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "true",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodPut,
		hasQuery: map[string][]string{
			"dc": {"dc1"},
		},
		hasBody: "someValue",
	})
	defer ts.Close()

	err := client.Put(ctx, "config/baz/bar", "someValue", Query{DC: "dc1"})
	require.NoError(t, err)
}

func Test_KV_Put_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodPut,
		hasQuery:  map[string][]string{},
		hasBody:   "someValue",
	})
	defer ts.Close()

	err := client.Put(ctx, "config/baz/bar", "someValue", Query{})
	require.EqualError(t, err, "status code (500)")
}

func Test_KV_Delete(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "true",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodDelete,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	err := client.Delete(ctx, "config/baz/bar", Query{})
	require.NoError(t, err)
}

func Test_KV_Delete_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      "true",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodDelete,
		hasQuery: map[string][]string{
			"dc": {"dc1"},
		},
	})
	defer ts.Close()

	err := client.Delete(ctx, "config/baz/bar", Query{DC: "dc1"})
	require.NoError(t, err)
}

func Test_KV_Delete_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/kv/config/baz/bar",
		hasMethod: http.MethodDelete,
		hasQuery:  map[string][]string{},
	})
	defer ts.Close()

	err := client.Delete(ctx, "config/baz/bar", Query{})
	require.EqualError(t, err, "status code (500)")
}

func Test_KV_Keys(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_kv_config_baz-keys.json"),
		hasPath:   "/v1/kv/config/baz",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"keys": {"true"},
		},
	})
	defer ts.Close()

	values, err := client.Keys(ctx, "config/baz", Query{})
	require.NoError(t, err)
	require.Equal(t, 3, len(values))
}

func Test_KV_Keys_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_kv_config_baz-keys-dc.json"),
		hasPath:   "/v1/kv/config/baz",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"keys": {"true"},
			"dc":   {"dc1"},
		},
	})
	defer ts.Close()

	values, err := client.Keys(ctx, "config/baz", Query{DC: "dc1"})
	require.NoError(t, err)
	require.Equal(t, 3, len(values))
	require.Equal(t, []string{
		"config/aaa/", "config/aaa/bbb", "config/aaa/ccc",
	}, values)
}

func Test_KV_Keys_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/kv/config/baz",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"keys": {"true"},
		},
	})
	defer ts.Close()

	_, err := client.Keys(ctx, "config/baz", Query{})
	require.EqualError(t, err, "status code (500)")
}

func Test_KV_Recurse(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_kv_config_baz-recurse.json"),
		hasPath:   "/v1/kv/config/baz",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"recurse": {"true"},
		},
	})
	defer ts.Close()

	values, err := client.Recurse(ctx, "config/baz", Query{})
	require.NoError(t, err)
	require.Equal(t, 8, len(values))
	require.Equal(t, []Pair{
		{Key: "config/baz", Value: ""},
		{Key: "config/baz/", Value: ""},
		{Key: "config/baz/bar", Value: "myValue"},
		{Key: "config/baz/cat", Value: "kitten"},
		{Key: "config/baz/sub/", Value: ""},
		{Key: "config/baz/sub/more/", Value: ""},
		{Key: "config/baz/sub/more/aaa", Value: "bbb"},
		{Key: "config/baz/sub/www", Value: "sdf"},
	}, values)
}

func Test_KV_Recurse_dc(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusOK,
		body:      load(t, "v1_kv_config_baz-recurse.json"),
		hasPath:   "/v1/kv/config/baz",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"recurse": {"true"},
			"dc":      {"dc1"},
		},
	})
	defer ts.Close()

	values, err := client.Recurse(ctx, "config/baz", Query{DC: "dc1"})
	require.NoError(t, err)
	require.Equal(t, 8, len(values))
	require.Equal(t, []Pair{
		{Key: "config/baz", Value: ""},
		{Key: "config/baz/", Value: ""},
		{Key: "config/baz/bar", Value: "myValue"},
		{Key: "config/baz/cat", Value: "kitten"},
		{Key: "config/baz/sub/", Value: ""},
		{Key: "config/baz/sub/more/", Value: ""},
		{Key: "config/baz/sub/more/aaa", Value: "bbb"},
		{Key: "config/baz/sub/www", Value: "sdf"},
	}, values)
}

func Test_KV_Recurse_err(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusInternalServerError,
		body:      "malfunction",
		hasPath:   "/v1/kv/config/baz",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"recurse": {"true"},
		},
	})
	defer ts.Close()

	_, err := client.Recurse(ctx, "config/baz", Query{})
	require.EqualError(t, err, "status code (500)")
}

func Test_KV_Recurse_non_existent(t *testing.T) {
	ctx, ts, client := testClient(&responder{
		t:         t,
		code:      http.StatusNotFound,
		body:      "malfunction",
		hasPath:   "/v1/kv/config/not-here",
		hasMethod: http.MethodGet,
		hasQuery: map[string][]string{
			"recurse": {"true"},
		},
	})
	defer ts.Close()

	_, err := client.Recurse(ctx, "config/not-here", Query{})
	require.EqualError(t, err, `key-space "config/not-here" does not exist`)
}
