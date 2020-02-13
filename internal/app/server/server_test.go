// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package server

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
)

var (
	natsURIs = flag.String("nats-uris", nats.DefaultURL, "Comma-separated list of NATS server URIs")
)

func TestNewRunStop(t *testing.T) {
	ctx := context.Background()

	nc, err := nats.Connect(*natsURIs)

	// Get a new server.
	c := Config{
		HTTPAddr: "localhost:",
		NATSConn: nc,
	}
	s, err := New(ctx, c)
	if err != nil {
		t.Fatalf("failed to get new server: %v", err)
	}

	wg := sync.WaitGroup{}

	// Start server goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Run()
	}()

	// Hit an endpoint to check the server is up.
	r, err := http.Get(fmt.Sprintf("http://%v/version", s.httpLn.Addr().String()))
	if err != nil {
		t.Errorf("failed to get HTTP: %v", err)
	}
	r.Body.Close()

	// Stop the server.
	if err := s.Stop(ctx); err != nil {
		t.Errorf("failed to stop server: %v", err)
	}

	// Wait until the server goroutine stops.
	wg.Wait()
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}
