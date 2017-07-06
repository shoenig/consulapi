// Author hoenig

package consulapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Client_Agent(t *testing.T) {
	client := New(ClientOptions{})
	defer cleanup(t, client)

	lanMembers, err := client.Members(false)
	require.NoError(t, err)
	require.Equal(t, 1, len(lanMembers))
	require.Equal(t, "dev-desktop1", lanMembers[0].Name)

	wanMembers, err := client.Members(true)
	require.NoError(t, err)
	require.Equal(t, 1, len(wanMembers))
	require.Equal(t, "dev-desktop1.dev", wanMembers[0].Name)

	err = client.Reload()
	require.NoError(t, err)

	err = client.MaintenanceMode(true, "for testing")
	require.NoError(t, err)

	err = client.MaintenanceMode(false, "")
	require.NoError(t, err)
}
