// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

// Persister is the interface by which all data is persisted.
type Persister interface {
	WorkflowPersister
}

// Resolver is the root type for resolving GraphQL queries.
type Resolver struct {
	p Persister
}

// New creates a new GraphQL resolver.
func New(p Persister) *Resolver {
	return &Resolver{p}
}
