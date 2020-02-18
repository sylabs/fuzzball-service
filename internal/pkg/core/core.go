// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

// Core represents core business logic.
type Core struct {
	p Persister
}

// Persister is the interface by which all data is persisted.
type Persister interface {
	WorkflowPersister
	JobPersister
	VolumePersister
}

// New creates a new core.
func New(p Persister) (*Core, error) {
	return &Core{
		p: p,
	}, nil
}
