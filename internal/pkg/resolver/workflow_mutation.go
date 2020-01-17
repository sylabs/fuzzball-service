// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// CreateWorkflow creates a new workflow.
func (r Resolver) CreateWorkflow(ctx context.Context, args struct {
	Name string
}) (*WorkflowResolver, error) {
	j := model.Workflow{
		Name: args.Name,
	}
	j, err := r.p.CreateWorkflow(ctx, j)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{j}, nil
}

// DeleteWorkflow deletes a workflow.
func (r Resolver) DeleteWorkflow(ctx context.Context, args struct {
	ID string
}) (*WorkflowResolver, error) {
	j, err := r.p.DeleteWorkflow(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{j}, nil
}
