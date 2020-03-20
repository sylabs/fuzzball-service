// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

// Servicer is the interface required to service GraphQL queries.
type Servicer interface {
	BuildInfoServicer
	UserServicer
	WorkflowServicer
}

// Resolver is the root type for resolving GraphQL queries.
type Resolver struct {
	s Servicer
}

// New creates a new GraphQL resolver.
func New(s Servicer) (*Resolver, error) {
	return &Resolver{s: s}, nil
}
