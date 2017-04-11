package consulapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Client(t *testing.T) {
	client := New(ClientOptions{})

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
}

func Test_fixup(t *testing.T) {
	tests := []struct {
		prefix string
		path   string
		params []string
		exp    string
	}{
		{prefix: "/v1/kv", path: "/devops/k1", exp: "/v1/kv/devops/k1"},
	}

	for _, test := range tests {
		fixed := fixup(test.prefix, test.path, test.params...)
		t.Log("fixed", fixed)
		require.Equal(t, test.exp, fixed)
	}
}
