// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"os"
	"runtime"

	"github.com/pbnjay/memory"
)

// Servicer is the interface required to service GraphQL queries.
type Servicer interface {
	UserServicer
	WorkflowServicer
}

// Resolver is the root type for resolving GraphQL queries.
type Resolver struct {
	s  Servicer
	si SystemInfo
}

// New creates a new GraphQL resolver.
func New(s Servicer) (*Resolver, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Resolver{
		s: s,
		si: SystemInfo{
			HostName:        hostName,
			CPUArchitecture: runtime.GOARCH,
			OSPlatform:      runtime.GOOS,
			Memory:          memory.TotalMemory() / 1024 / 1024, // megabytes
		},
	}, nil
}
