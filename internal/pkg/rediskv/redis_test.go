// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package rediskv

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
)

var (
	redisURI = flag.String("redis-uri", "localhost:6379", "URI of Redis")

	testConnection *Connection
)

func run(m *testing.M) int {
	// Create test connection.
	c, err := NewConnection(*redisURI)
	if err != nil {
		log.Printf("failed to create new connection: %v", err)
		return -1
	}
	defer c.Disconnect()
	testConnection = c

	return m.Run()
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(run(m))
}

func TestGetSet(t *testing.T) {
	i := rand.Int()
	testkey, testvalue := fmt.Sprintf("testkey-%d", i), fmt.Sprintf("testvalue-%d", i)
	_, err := testConnection.Get(testkey)
	if err == nil {
		t.Fatal("unexpected success")
	}

	err = testConnection.Set(testkey, testvalue)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	val, err := testConnection.Get(testkey)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	if testvalue != val {
		t.Fatalf("want %s, got %s", string(testvalue), string(val))
	}
}
