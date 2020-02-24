// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/core"
)

// WorkflowServicer is the interface by which workflows are serviced.
type WorkflowServicer interface {
	CreateWorkflow(context.Context, core.WorkflowSpec) (core.Workflow, error)
	DeleteWorkflow(context.Context, string) (core.Workflow, error)
	GetWorkflow(context.Context, string) (core.Workflow, error)
}

// WorkflowResolver resolves a workflow.
type WorkflowResolver struct {
	w core.Workflow
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
func (r *WorkflowResolver) CreatedBy(ctx context.Context) (*UserResolver, error) {
	u, err := r.w.CreatedBy(ctx)
	if err != nil {
		return nil, err
	}
	return &UserResolver{u: &u}, nil
}

// CreatedAt resolves when the workflow was created.
func (r *WorkflowResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.w.CreatedAt}
}

// StartedAt returns when the workflow started, if it has started.
func (r *WorkflowResolver) StartedAt() *graphql.Time {
	if t := r.w.StartedAt; t != nil {
		return &graphql.Time{Time: *t}
	}
	return nil
}

// FinishedAt returns when the workflow finished, if it has finished.
func (r *WorkflowResolver) FinishedAt() *graphql.Time {
	if t := r.w.FinishedAt; t != nil {
		return &graphql.Time{Time: *t}
	}
	return nil
}

// Status resolves the state of the workflow.
func (r *WorkflowResolver) Status() string {
	return r.w.Status
}

// Jobs looks up jobs associated with the workflow.
func (r *WorkflowResolver) Jobs(ctx context.Context, args pageArgs) (*JobConnectionResolver, error) {
	p, err := r.w.JobsPage(ctx, convertPageArgs(args))
	if err != nil {
		return nil, err
	}
	return &JobConnectionResolver{p}, nil
}

// Volumes looks up volumes associated with the workflow.
func (r *WorkflowResolver) Volumes(ctx context.Context, args pageArgs) (*VolumeConnectionResolver, error) {
	p, err := r.w.VolumesPage(ctx, convertPageArgs(args))
	if err != nil {
		return nil, err
	}
	return &VolumeConnectionResolver{p}, nil
}
