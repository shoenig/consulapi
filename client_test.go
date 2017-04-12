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
