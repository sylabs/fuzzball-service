// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build integration

package mongodb

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var (
	mongoURI = flag.String("mongo-uri", "mongodb://localhost", "URI of MongoDB database")

	testDBName     string
	testConnection *Connection
)

func TestNewDisconnect(t *testing.T) {
	ctx := context.Background()
	expiredCtx, cancel := context.WithDeadline(ctx, time.Now().Add(-time.Hour))
	defer cancel()

	tests := []struct {
		name    string
		ctx     context.Context
		uri     string
		wantErr bool
	}{
		{"Success", ctx, *mongoURI, false},
		{"ExpiredContext", expiredCtx, *mongoURI, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewConnection(tt.ctx, tt.uri, testDBName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if err := c.Disconnect(tt.ctx); err != nil {
					t.Fatalf("failed to disconnect: %v", err)
				}
			}
		})
	}
}

func run(m *testing.M) int {
	ctx := context.Background()

	// Construct unique database name.
	testDBName = fmt.Sprintf("compute-service-%09d", time.Now().UnixNano()%time.Second.Nanoseconds())

	// Create test connection.
	c, err := NewConnection(ctx, *mongoURI, testDBName)
	if err != nil {
		log.Printf("failed to create new connection: %v", err)
		return -1
	}
	defer c.Disconnect(ctx)
	testConnection = c

	return m.Run()
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(run(m))
}
