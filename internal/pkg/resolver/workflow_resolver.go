// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/core"
)

// WorkflowServicer is the interface by which workflows are serviced.
type WorkflowServicer interface {
	CreateWorkflow(context.Context, core.WorkflowSpec) (core.Workflow, error)
	DeleteWorkflow(context.Context, string) (core.Workflow, error)
	GetWorkflow(context.Context, string) (core.Workflow, error)
	GetWorkflows(context.Context, core.PageArgs) (core.WorkflowsPage, error)
}

// WorkflowResolver resolves a workflow.
type WorkflowResolver struct {
	w core.Workflow
	s Servicer
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
		u: &core.User{
			ID:    "507f1f77bcf86cd799439011",
			Login: "jimbob",
		},
		s: r.s,
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

// Status resolves the state of the workflow.
func (r *WorkflowResolver) Status() string {
	return r.w.Status
}

// Jobs looks up jobs associated with the workflow.
func (r *WorkflowResolver) Jobs(ctx context.Context, args struct {
	After  *string
	Before *string
	First  *int
	Last   *int
}) (*JobConnectionResolver, error) {
	pa := core.PageArgs{
		After:  args.After,
		Before: args.Before,
		First:  args.First,
		Last:   args.Last,
	}
	p, err := r.s.GetJobsByWorkflowID(ctx, pa, r.w.ID)
	if err != nil {
		return nil, err
	}
	return &JobConnectionResolver{p, r.s}, nil
}

// Volumes looks up volumes associated with the workflow.
func (r *WorkflowResolver) Volumes(ctx context.Context, args struct {
	After  *string
	Before *string
	First  *int
	Last   *int
}) (*VolumeConnectionResolver, error) {
	pa := core.PageArgs{
		After:  args.After,
		Before: args.Before,
		First:  args.First,
		Last:   args.Last,
	}
	p, err := r.s.GetVolumesByWorkflowID(ctx, pa, r.w.ID)
	if err != nil {
		return nil, err
	}
	return &VolumeConnectionResolver{p, r.s}, nil
}
