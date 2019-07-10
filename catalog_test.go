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

	services, err := client.Services("")
	require.NoError(t, err)
	require.Equal(t, 3, len(services))
	expServices := map[string][]string{
		"consul":   {},
		"myapp":    {"myapp-tag1", "myapp-tag2"},
		"otherapp": {"otherapp-tag1", "otherapp-tag2", "otherapp-tag3"},
	}
	require.Equal(t, expServices, services)

	myapp, err := client.Service("dev", "myapp")
	require.NoError(t, err)
	require.Equal(t, 1, len(myapp))
	require.Equal(t, "myapp", myapp[0].ServiceID)
	require.Equal(t, 29019, myapp[0].ServicePort)

	otherapp, err := client.Service("dev", "otherapp", "otherapp-tag3")
	require.NoError(t, err)
	require.Equal(t, 1, len(otherapp))
	require.Equal(t, "otherapp", otherapp[0].ServiceName)
	require.Equal(t, "otherapp2", otherapp[0].ServiceID)
	require.Equal(t, 3, len(otherapp[0].ServiceTags))
	require.Equal(t, "otherapp-tag3", otherapp[0].ServiceTags[2])

	nodeinfo, err := client.Node("dev", "dev-desktop1")
	require.NoError(t, err)

	require.Equal(t, "dev-desktop1", nodeinfo.Node.Name)
	require.Equal(t, 4, len(nodeinfo.Services))
}
