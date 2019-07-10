package consulapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

	_, err = client.Get("", "/test/random_path_to_force_404")
	require.NotNil(t, err)
	requestError, ok := err.(*RequestError)
	require.True(t, ok, "Error must be RequestError")
	require.Equal(t, 404, requestError.StatusCode())
}
