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
	SetJobStatus(context.Context, string, string, int) error
}

// Scheduler represents an instance of the scheduler.
type Scheduler struct {
	m Messager
	p Persister
}

// New creates a new scheduler.
func New(m Messager, p Persister) (*Scheduler, error) {
	return &Scheduler{m, p}, nil
}
