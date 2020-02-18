// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// UserResolver resolves a user.
type UserResolver struct {
	u *model.User
	p Persister
	f IOFetcher
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
func (r *UserResolver) Workflows(ctx context.Context, args struct {
	After  *string
	Before *string
	First  *int32
	Last   *int32
}) (*WorkflowConnectionResolver, error) {
	pa := model.PageArgs{
		After:  args.After,
		Before: args.Before,
	}
	if args.First != nil {
		first := int(*args.First)
		pa.First = &first
	}
	if args.Last != nil {
		last := int(*args.Last)
		pa.Last = &last
	}
	p, err := r.p.GetWorkflows(ctx, pa)
	if err != nil {
		return nil, err
	}
	return &WorkflowConnectionResolver{p, r.p, r.f}, nil
}

// Jobs looks up jobs associated with the user.
func (r *UserResolver) Jobs(ctx context.Context, args struct {
	After  *string
	Before *string
	First  *int32
	Last   *int32
}) (*JobConnectionResolver, error) {
	pa := model.PageArgs{
		After:  args.After,
		Before: args.Before,
	}
	if args.First != nil {
		first := int(*args.First)
		pa.First = &first
	}
	if args.Last != nil {
		last := int(*args.Last)
		pa.Last = &last
	}
	p, err := r.p.GetJobs(ctx, pa)
	if err != nil {
		return nil, err
	}
	return &JobConnectionResolver{p, r.p, r.f}, nil
}

// Volumes looks up volumes associated with the user.
func (r *UserResolver) Volumes(ctx context.Context, args struct {
	After  *string
	Before *string
	First  *int32
	Last   *int32
}) (*VolumeConnectionResolver, error) {
	pa := model.PageArgs{
		After:  args.After,
		Before: args.Before,
	}
	if args.First != nil {
		first := int(*args.First)
		pa.First = &first
	}
	if args.Last != nil {
		last := int(*args.Last)
		pa.Last = &last
	}
	p, err := r.p.GetVolumes(ctx, pa)
	if err != nil {
		return nil, err
	}
	return &VolumeConnectionResolver{p, r.p}, nil
}
