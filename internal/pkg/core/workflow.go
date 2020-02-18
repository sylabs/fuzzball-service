// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// WorkflowPersister is the interface by which workflows are persisted.
type WorkflowPersister interface {
	CreateWorkflow(context.Context, model.Workflow) (model.Workflow, error)
	DeleteWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflows(context.Context, model.PageArgs) (model.WorkflowsPage, error)
}

// CreateWorkflow creates a new workflow. If an ID is provided in w, it is ignored and replaced
// with a unique identifier in the returned workflow.
func (c *Core) CreateWorkflow(ctx context.Context, w model.Workflow) (model.Workflow, error) {
	return c.p.CreateWorkflow(ctx, w)
}

// DeleteWorkflow deletes a workflow by ID. If the supplied ID is not valid, or there there is not
// a workflow with a matching ID in the database, an error is returned.
func (c *Core) DeleteWorkflow(ctx context.Context, id string) (w model.Workflow, err error) {
	return c.p.DeleteWorkflow(ctx, id)
}

// GetWorkflow retrieves a workflow by ID. If the supplied ID is not valid, or there there is not a
// workflow with a matching ID in the database, an error is returned.
func (c *Core) GetWorkflow(ctx context.Context, id string) (w model.Workflow, err error) {
	return c.p.GetWorkflow(ctx, id)
}

// GetWorkflows returns a list of all workflows.
func (c *Core) GetWorkflows(ctx context.Context, pa model.PageArgs) (p model.WorkflowsPage, err error) {
	return c.p.GetWorkflows(ctx, pa)
}
