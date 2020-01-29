// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import "context"

// Persister is the interface that describes what is needed to persist scheduler data.
type Persister interface {
	SetWorkflowStatus(context.Context, string, string) error
	SetJobStatus(context.Context, string, string, int) error
}

// Scheduler represents an instance of the scheduler.
type Scheduler struct {
	p Persister
}

// New creates a new scheduler.
func New(p Persister) (*Scheduler, error) {
	return &Scheduler{p}, nil
}
