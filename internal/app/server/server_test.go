// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestNewRunStop(t *testing.T) {
	ctx := context.Background()

	// Mock key server.
	jwks := mockJWKS{keys: testKeySet}
	mks := httptest.NewServer(&jwks)
	defer mks.Close()

	// Mock discovery server.
	m := mockOAuthDisco{md: testMetadata}
	m.md.JWKSURI = mks.URL
	mds := httptest.NewServer(&m)
	m.md.Issuer = mds.URL
	defer mds.Close()

	// Get a new server.
	c := Config{
		HTTPAddr:        "localhost:",
		OAuth2IssuerURI: mds.URL,
	}
	s, err := New(ctx, nil, c)
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
	r, err := http.Get(fmt.Sprintf("http://%v/graphql", s.httpLn.Addr().String()))
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
