// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// UserResolver resolves a user.
type UserResolver struct {
	u  *model.User
	wp WorkflowPersister
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
	First  *int
	Last   *int
}) (*WorkflowConnectionResolver, error) {
	pa := model.PageArgs{
		After:  args.After,
		Before: args.Before,
		First:  args.First,
		Last:   args.Last,
	}
	p, err := r.wp.GetWorkflows(ctx, pa)
	if err != nil {
		return nil, err
	}
	return &WorkflowConnectionResolver{p}, nil
}
