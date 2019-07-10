package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/shoenig/consulapi"
)

var name string

// Sample code demonstrating a leader electing client which does "work"
// inside a consulapi.AsLeaderFunc, respecting the provided context.Context
// by returning when the context is marked as Done.

func main() {
	name = os.Args[1]
	log.Println("starting elector with name: " + name)

	consul := consulapi.New(consulapi.ClientOptions{
		Logger: log.New(os.Stdout, "<consulapi> ", log.LstdFlags),
	})

	leadershipConfig := consulapi.LeadershipConfig{
		Key:         "service/elector/leader",
		Description: "elector-leader-session",
		LockDelay:   3 * time.Second,
		TTL:         10 * time.Second,
		ContactInfo: fmt.Sprintf("elector:%s", name),
	}

	f := func(ctx context.Context) error {
		log.Printf("--- leaderfunc called ---")

		select {
		case <-ctx.Done():
			log.Printf("--- leaderfunc context was signaled as done ---")
			return errors.Wrap(ctx.Err(), "context is done")
		}
	}

	session, err := consul.Participate(leadershipConfig, f)
	if err != nil {
		panic(err)
	}

	log.Printf("[elector %s] going to idle, with session id: %s", name, session.SessionID())

	for range time.Tick(2 * time.Second) {
		showLeader(session)
	}
}

func showLeader(session consulapi.LeaderSession) {
	leader, err := session.Current()
	if err != nil {
		log.Printf("[elector %s] could not get current leader: %v", name, err)
	} else {
		log.Printf("[elector %s] looked up current leader which is: %s", name, leader)
	}
}
