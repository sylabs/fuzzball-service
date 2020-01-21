// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// WorkflowPersister is the interface by which workflows are persisted.
type WorkflowPersister interface {
	CreateWorkflow(context.Context, model.Workflow) (model.Workflow, error)
	DeleteWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflows(context.Context, model.PageArgs) (model.WorkflowsPage, error)
}

// WorkflowResolver resolves a workflow.
type WorkflowResolver struct {
	w model.Workflow
	p Persister
}

// ID resolves the workflow ID.
func (r *WorkflowResolver) ID() graphql.ID {
	return graphql.ID(r.w.ID)
}

// Name resolves the workflow name.
func (r *WorkflowResolver) Name() string {
	return r.w.Name
}

// CreatedBy resolves the user who created the workflow.
func (r *WorkflowResolver) CreatedBy() *UserResolver {
	return &UserResolver{
		u: &model.User{
			ID:    "507f1f77bcf86cd799439011",
			Login: "jimbob",
		},
		p: r.p,
	}
}

// CreatedAt resolves when the workflow was created.
func (r *WorkflowResolver) CreatedAt() (graphql.Time, error) {
	return graphql.Time{Time: time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC)}, nil // TODO
}

// StartedAt returns when the workflow started, if it has started.
func (r *WorkflowResolver) StartedAt() *graphql.Time {
	return nil // TODO
}

// FinishedAt returns when the workflow finished, if it has finished.
func (r *WorkflowResolver) FinishedAt() *graphql.Time {
	return nil // TODO
}
