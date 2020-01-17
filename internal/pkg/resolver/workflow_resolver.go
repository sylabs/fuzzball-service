// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// WorkflowPersister is the interface by which workflows are persisted.
type WorkflowPersister interface {
	CreateWorkflow(context.Context, model.Workflow) (model.Workflow, error)
	DeleteWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflows(context.Context) (model.WorkflowsPage, error)
}

// WorkflowResolver resolves a workflow.
type WorkflowResolver struct {
	w *model.Workflow
}

// ID resolves the workflow ID.
func (r *WorkflowResolver) ID() graphql.ID {
	return graphql.ID(r.w.ID)
}

// Name resolves the workflow name.
func (r *WorkflowResolver) Name() string {
	return r.w.Name
}
