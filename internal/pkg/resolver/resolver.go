// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"os"
	"runtime"

	"github.com/pbnjay/memory"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// IOFetcher is the interface where IO data is retrieved.
type IOFetcher interface {
	JobOutputFetcher
}

// Persister is the interface by which all data is persisted.
type Persister interface {
	WorkflowPersister
	JobPersister
	VolumePersister
}

// Scheduler is the interface by which all workflows are scheduled.
type Scheduler interface {
	AddWorkflow(context.Context, model.Workflow, []model.Job, map[string]model.Volume) error
}

// Resolver is the root type for resolving GraphQL queries.
type Resolver struct {
	p  Persister
	f  IOFetcher
	s  Scheduler
	si SystemInfo
}

// New creates a new GraphQL resolver.
func New(p Persister, f IOFetcher, s Scheduler) (*Resolver, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Resolver{
		p: p,
		f: f,
		s: s,
		si: SystemInfo{
			HostName:        hostName,
			CPUArchitecture: runtime.GOARCH,
			OSPlatform:      runtime.GOOS,
			Memory:          memory.TotalMemory() / 1024 / 1024, // megabytes
		},
	}, nil
}
