package consulapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Running these tests requires running the consul defined in the hack/docker-compose.yml file.
// With docker and docker-compose installed, simply run "docker-compose up" in the hack directory
// to get up and running.

func cleanup(t *testing.T, client Client) {
	t.Log("-- cleaning up consul key space --")
	keys, err := client.Keys("", "/")
	require.NoError(t, err)

	for _, key := range keys {
		err := client.Delete("", key)
		require.NoError(t, err)
	}
}

func Test_Client_KV(t *testing.T) {
	client := New(ClientOptions{})
	defer cleanup(t, client)

	dcs, err := client.Datacenters()
	require.NoError(t, err)
	require.Equal(t, 1, len(dcs))
	require.Equal(t, "dev", dcs[0])

	err = client.Put("", "/test/t1", "value1")
	require.NoError(t, err)

	value, err := client.Get("", "/test/t1")
	require.NoError(t, err)
	require.Equal(t, "value1", value)

	err = client.Delete("", "/test/t1")
	require.NoError(t, err)

	_, err = client.Get("", "/test/t1")
	require.Error(t, err)

	err = client.Put("dev", "/test/t2", "value2")
	require.NoError(t, err)
	err = client.Put("dev", "/other/t3", "value3")
	require.NoError(t, err)

	keys, err := client.Keys("", "/")
	require.NoError(t, err)
	require.Equal(t, 2, len(keys))
	require.Equal(t, "other/t3", keys[0])
	require.Equal(t, "test/t2", keys[1])

	err = client.Put("", "other/t4", "value4")
	require.NoError(t, err)
	err = client.Put("", "other/t4/sub1", "value5")

	all, err := client.Recurse("dev", "/other")
	require.NoError(t, err)
	require.Equal(t, 3, len(all))

	require.Equal(t, "other/t3", all[0][0])
	require.Equal(t, "value3", all[0][1])
}

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
