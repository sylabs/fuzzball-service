// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
)

// Persister is the interface by which all data is persisted.
type Persister interface {
	WorkflowPersister
	JobPersister
	VolumePersister
}

// IOFetcher is the interface where IO data is retrieved.
type IOFetcher interface {
	JobOutputFetcher
}

// Scheduler is the interface by which all workflows are scheduled.
type Scheduler interface {
	AddWorkflow(context.Context, Workflow, []Job, map[string]Volume) error
}

// Core represents core business logic.
type Core struct {
	p Persister
	f IOFetcher
	s Scheduler
}

// New creates a new core.
func New(p Persister, f IOFetcher, s Scheduler) (*Core, error) {
	return &Core{p: p, f: f, s: s}, nil
}
