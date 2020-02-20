// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
)

// Job contains information about an indivisual job.
type Job struct {
	ID         string              `bson:"_id,omitempty"`
	WorkflowID string              `bson:"workflowID"`
	Name       string              `bson:"name"`
	Image      string              `bson:"image"`
	Command    []string            `bson:"command"`
	Status     string              `bson:"status"`
	ExitCode   int                 `bson:"exitCode"`
	Requires   []string            `bson:"requires"`
	Volumes    []VolumeRequirement `bson:"volumes"`

	c *Core // Used internally for lazy loading.
}

// VolumeRequirement describes a required volume.
type VolumeRequirement struct {
	VolumeID string `bson:"volumeID"`
	Name     string `bson:"name"`
	Location string `bson:"location"`
}

// WithCore returns a job with the core field set to c.
func (j Job) WithCore(c *Core) Job {
	j.c = c
	return j
}

// GetOutput retrieves the output of the job.
func (j Job) GetOutput() (string, error) {
	return j.c.f.GetJobOutput(j.ID)
}

// JobsPage represents a page of jobs resulting from a query, and associated metadata.
type JobsPage struct {
	Jobs       []Job    // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
}

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	CreateJob(context.Context, Job) (Job, error)
	DeleteJobsByWorkflowID(context.Context, string) error
	GetJob(context.Context, string) (Job, error)
	GetJobs(context.Context, PageArgs) (JobsPage, error)
	GetJobsByWorkflowID(context.Context, PageArgs, string) (JobsPage, error)
	GetJobsByID(context.Context, PageArgs, string, []string) (JobsPage, error)
}

// CreateJob creates a new job. If an ID is provided in j, it is ignored and replaced with a unique
// identifier in the returned job.
func (c *Core) CreateJob(ctx context.Context, j Job) (Job, error) {
	j, err := c.p.CreateJob(ctx, j)
	return j.WithCore(c), err
}

// DeleteJobsByWorkflowID deletes jobs with the given workflow ID.
func (c *Core) DeleteJobsByWorkflowID(ctx context.Context, wid string) error {
	return c.p.DeleteJobsByWorkflowID(ctx, wid)
}

// GetJob retrieves a job by ID. If the supplied ID is not valid, or there there is not a
// job with a matching ID in the database, an error is returned.
func (c *Core) GetJob(ctx context.Context, id string) (Job, error) {
	j, err := c.p.GetJob(ctx, id)
	return j.WithCore(c), err
}

// GetJobs returns a list of all jobs.
func (c *Core) GetJobs(ctx context.Context, pa PageArgs) (JobsPage, error) {
	p, err := c.p.GetJobs(ctx, pa)
	for i, j := range p.Jobs {
		p.Jobs[i] = j.WithCore(c)
	}
	return p, err
}

// GetJobsByWorkflowID returns a list of all jobs for a given workflow.
func (c *Core) GetJobsByWorkflowID(ctx context.Context, pa PageArgs, wid string) (JobsPage, error) {
	p, err := c.p.GetJobsByWorkflowID(ctx, pa, wid)
	for i, j := range p.Jobs {
		p.Jobs[i] = j.WithCore(c)
	}
	return p, err
}

// GetJobsByID returns a list of jobs by name within a given workflow.
func (c *Core) GetJobsByID(ctx context.Context, pa PageArgs, wid string, ids []string) (JobsPage, error) {
	p, err := c.p.GetJobsByID(ctx, pa, wid, ids)
	for i, j := range p.Jobs {
		p.Jobs[i] = j.WithCore(c)
	}
	return p, err
}
