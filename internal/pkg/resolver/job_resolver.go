// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// JobPersister is the interface by which workflows are persisted.
type JobPersister interface {
	CreateJob(context.Context, model.Job) (model.Job, error)
	DeleteJobsByWorkflowID(context.Context, string) error
	GetJob(context.Context, string) (model.Job, error)
	GetJobs(context.Context, model.PageArgs) (model.JobsPage, error)
	GetJobsByWorkflowID(context.Context, model.PageArgs, string) (model.JobsPage, error)
}

// JobResolver resolves a workflow.
type JobResolver struct {
	j model.Job
	p Persister
}

// ID resolves the Job ID.
func (r *JobResolver) ID() graphql.ID {
	return graphql.ID(r.j.ID)
}

// Name resolves the job name.
func (r *JobResolver) Name() string {
	return r.j.Name
}

// Image resolves the job image.
func (r *JobResolver) Image() string {
	return r.j.Image
}

// Command resolves the job command.
func (r *JobResolver) Command() []string {
	return r.j.Command
}

// CreatedBy resolves the user who created the job.
func (r *JobResolver) CreatedBy() *UserResolver {
	return &UserResolver{
		u: &model.User{
			ID:    "507f1f77bcf86cd799439011",
			Login: "jimbob",
		},
		p: r.p,
	}
}

// CreatedAt resolves when the job was created.
func (r *JobResolver) CreatedAt() (graphql.Time, error) {
	return graphql.Time{Time: time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC)}, nil // TODO
}

// StartedAt returns when the job started, if it has started.
func (r *JobResolver) StartedAt() *graphql.Time {
	return nil // TODO
}

// FinishedAt returns when the job finished, if it has finished.
func (r *JobResolver) FinishedAt() *graphql.Time {
	return nil // TODO
}
