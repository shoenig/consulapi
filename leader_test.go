package consulapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Client_LeadershipManager(t *testing.T) {
	client := New(ClientOptions{})
	defer cleanup(t, client)

	f := func(context.Context) error {
		fmt.Println("f was called")
		return nil
	}

	t.Log("starting leadership manager")
	manager, err := client.Participate(LeadershipConfig{
		Key:         "service/test/leader",
		Description: "my-leadership-session",
		LockDelay:   5 * time.Second,
		TTL:         10 * time.Second,
	}, f)
	require.NoError(t, err)

	t.Log("created client")

	time.Sleep(5 * time.Second)
	t.Log("slept, abdicating")

	err = manager.Abdicate()
	require.NoError(t, err)

	t.Log("exiting")
}

func Test_Client_LeadershipManager_3(t *testing.T) {
	clients := []Client{
		New(ClientOptions{}),
		New(ClientOptions{}),
		New(ClientOptions{}),
	}

	defer cleanup(t, clients[0])
	defer cleanup(t, clients[1])
	defer cleanup(t, clients[2])

	f := func(context.Context) error {
		fmt.Println("leadership starting")
		time.Sleep(1 * time.Second)
		fmt.Println("leadership exiting")
		return nil
	}

	lc := LeadershipConfig{
		Key:         "service/test3/leader",
		Description: "my-3-leadership-session",
		LockDelay:   1 * time.Second,
		TTL:         10 * time.Second,
	}

	lm1, err := clients[0].Participate(lc, f)
	require.NoError(t, err)
	lm2, err := clients[1].Participate(lc, f)
	require.NoError(t, err)
	lm3, err := clients[2].Participate(lc, f)
	require.NoError(t, err)

	time.Sleep(2 * time.Second)

	_ = lm1.Abdicate()
	_ = lm2.Abdicate()
	_ = lm3.Abdicate()

	time.Sleep(4 * time.Second)
}
