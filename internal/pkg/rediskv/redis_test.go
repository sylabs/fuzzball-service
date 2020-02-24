// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package rediskv

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"testing"
)

var (
	redisURI = flag.String("redis-uri", "redis://localhost", "URI of Redis")

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

func TestGetSetAppend(t *testing.T) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt32)))
	if err != nil {
		t.Fatalf("failed to generate random int: %v", err)
	}
	i := int32(n.Int64())
	testkey, testvalue := fmt.Sprintf("testkey-%d", i), fmt.Sprintf("testvalue-%d", i)
	val, err := testConnection.Get(testkey)
	if err != nil {
		t.Fatal("unexpected failure")
	}
	if val != "" {
		t.Fatalf("want %q, got %q", "", val)
	}

	err = testConnection.Set(testkey, testvalue)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	val, err = testConnection.Get(testkey)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	if testvalue != val {
		t.Fatalf("want %s, got %s", testvalue, val)
	}

	err = testConnection.Append(testkey, testvalue)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	val, err = testConnection.Get(testkey)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	want := testvalue + testvalue
	if want != val {
		t.Fatalf("want %q, got %q", want, val)
	}
}
