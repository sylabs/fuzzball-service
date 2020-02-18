// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	CreateJob(context.Context, model.Job) (model.Job, error)
	DeleteJobsByWorkflowID(context.Context, string) error
	GetJob(context.Context, string) (model.Job, error)
	GetJobs(context.Context, model.PageArgs) (model.JobsPage, error)
	GetJobsByWorkflowID(context.Context, model.PageArgs, string) (model.JobsPage, error)
	GetJobsByID(context.Context, model.PageArgs, string, []string) (model.JobsPage, error)
}

// CreateJob creates a new job. If an ID is provided in j, it is ignored and replaced with a unique
// identifier in the returned job.
func (c *Core) CreateJob(ctx context.Context, j model.Job) (model.Job, error) {
	return c.p.CreateJob(ctx, j)
}

// DeleteJobsByWorkflowID deletes jobs with the given workflow ID.
func (c *Core) DeleteJobsByWorkflowID(ctx context.Context, wid string) error {
	return c.p.DeleteJobsByWorkflowID(ctx, wid)
}

// GetJob retrieves a job by ID. If the supplied ID is not valid, or there there is not a
// job with a matching ID in the database, an error is returned.
func (c *Core) GetJob(ctx context.Context, id string) (j model.Job, err error) {
	return c.p.GetJob(ctx, id)
}

// GetJobs returns a list of all jobs.
func (c *Core) GetJobs(ctx context.Context, pa model.PageArgs) (p model.JobsPage, err error) {
	return c.p.GetJobs(ctx, pa)
}

// GetJobsByWorkflowID returns a list of all jobs for a given workflow.
func (c *Core) GetJobsByWorkflowID(ctx context.Context, pa model.PageArgs, wid string) (p model.JobsPage, err error) {
	return c.p.GetJobsByWorkflowID(ctx, pa, wid)
}

// GetJobsByID returns a list of jobs by name within a given workflow.
func (c *Core) GetJobsByID(ctx context.Context, pa model.PageArgs, wid string, ids []string) (p model.JobsPage, err error) {
	return c.p.GetJobsByID(ctx, pa, wid, ids)
}
