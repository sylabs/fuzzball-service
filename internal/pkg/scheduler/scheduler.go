// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

// Messager is the interface that is needed to send and receive messages.
type Messager interface {
	Request(subject string, v interface{}, vPtr interface{}, timeout time.Duration) error
	Subscribe(subject string, cb nats.Handler) (*nats.Subscription, error)
}

// Persister is the interface that describes what is needed to persist scheduler data.
type Persister interface {
	SetWorkflowStatus(context.Context, string, string) error
	SetJobStatus(context.Context, string, string) error
	SetJobExitCode(context.Context, string, int) error
}

// IOPersister is the interface that describes what is needed to persist Job IO data.
type IOPersister interface {
	Set(string, string) error
	Get(string) (string, error)
}

// Scheduler represents an instance of the scheduler.
type Scheduler struct {
	m   Messager
	p   Persister
	iop IOPersister
}

// New creates a new scheduler.
func New(m Messager, p Persister, iop IOPersister) (*Scheduler, error) {
	return &Scheduler{m, p, iop}, nil
}
