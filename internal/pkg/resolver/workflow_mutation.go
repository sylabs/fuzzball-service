// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/core"
)

// CreateWorkflow creates a new workflow.
func (r Resolver) CreateWorkflow(ctx context.Context, args struct {
	Spec core.WorkflowSpec
}) (*WorkflowResolver, error) {
	w, err := r.p.CreateWorkflow(ctx, args.Spec)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{w, r.p}, nil
}

// DeleteWorkflow deletes a workflow.
func (r Resolver) DeleteWorkflow(ctx context.Context, args struct {
	ID string
}) (*WorkflowResolver, error) {
	w, err := r.p.DeleteWorkflow(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{w, r.p}, nil
}
