// Author hoenig

package consulapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Client_Catalog(t *testing.T) {
	client := New(ClientOptions{})
	defer cleanup(t, client)

	dcs, err := client.Datacenters()
	require.NoError(t, err)
	require.Equal(t, 1, len(dcs))
	require.Equal(t, "dev", dcs[0])

	nodes, err := client.Nodes("dev")
	require.NoError(t, err)
	require.Equal(t, 1, len(nodes))
	require.Equal(t, "dev-desktop1", nodes[0].Name)
}
