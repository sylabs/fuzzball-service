// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/core"
)

// UserServicer is the interface by which users are serviced.
type UserServicer interface {
	Viewer(ctx context.Context) (core.User, error)
}

// UserResolver resolves a user.
type UserResolver struct {
	u *core.User
}

// ID resolves the unique user ID.
func (r *UserResolver) ID() graphql.ID {
	return graphql.ID(r.u.ID)
}

// Login resolves the username used to login.
func (r *UserResolver) Login() string {
	return r.u.Login
}

// Workflows looks up workflows associated with the user.
func (r *UserResolver) Workflows(ctx context.Context, args pageArgs) (*WorkflowConnectionResolver, error) {
	p, err := r.u.WorkflowsPage(ctx, convertPageArgs(args))
	if err != nil {
		return nil, err
	}
	return &WorkflowConnectionResolver{p}, nil
}

// Jobs looks up jobs associated with the user.
func (r *UserResolver) Jobs(ctx context.Context, args pageArgs) (*JobConnectionResolver, error) {
	p, err := r.u.JobsPage(ctx, convertPageArgs(args))
	if err != nil {
		return nil, err
	}
	return &JobConnectionResolver{p}, nil
}

// Volumes looks up volumes associated with the user.
func (r *UserResolver) Volumes(ctx context.Context, args pageArgs) (*VolumeConnectionResolver, error) {
	p, err := r.u.VolumesPage(ctx, convertPageArgs(args))
	if err != nil {
		return nil, err
	}
	return &VolumeConnectionResolver{p}, nil
}
