// Author hoenig

package consulapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Client_Session(t *testing.T) {
	client := New(ClientOptions{})
	defer cleanup(t, client)

	id1, err := client.CreateSession("", SessionConfig{
		Node:      "dev-desktop1",
		Name:      "my_session1",
		TTL:       SessionMininumTTL,
		LockDelay: SessionMaximumLockDelay,
		Behavior:  SessionDelete,
	})
	require.NoError(t, err)
	t.Log("session id:", id1)

	id2, err := client.CreateSession("", SessionConfig{
		Node:      "dev-desktop1",
		Name:      "my_session2",
		TTL:       SessionMininumTTL,
		LockDelay: SessionMaximumLockDelay,
		Behavior:  SessionDelete,
	})
	require.NoError(t, err)
	t.Log("session id:", id2)

	sessions, err := client.ListSessions("", "dev-desktop1")
	require.NoError(t, err)
	t.Log("sessions:", sessions)

	session, err := client.ReadSession("", id1)
	require.NoError(t, err)
	t.Log("session1:", session)

	ttl, err := client.RenewSession("", id1)
	require.NoError(t, err)
	t.Log("new ttl:", ttl)

	err = client.DeleteSession("", id1)
	require.NoError(t, err)
}
