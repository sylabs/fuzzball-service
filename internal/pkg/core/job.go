// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
)

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	CreateJob(context.Context, Job) (Job, error)
	DeleteJobsByWorkflowID(context.Context, string) error
	GetJobs(context.Context, PageArgs) (JobsPage, error)
	GetJobsByWorkflowID(context.Context, PageArgs, string) (JobsPage, error)
	GetJobsByID(context.Context, PageArgs, string, []string) (JobsPage, error)
}

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

// setCore sets the core of j to c.
func (j *Job) setCore(c *Core) {
	j.c = c
}

// GetOutput retrieves the output of the job.
func (j Job) GetOutput() (string, error) {
	return j.c.f.GetJobOutput(j.ID)
}

// CreatedBy retrieves the user that created job j.
func (j Job) CreatedBy(ctx context.Context) (User, error) {
	u := User{
		ID:    "507f1f77bcf86cd799439011",
		Login: "jimbob",
	}
	u.setCore(j.c)
	return u, nil
}

// JobsPage represents a page of jobs resulting from a query, and associated metadata.
type JobsPage struct {
	Jobs       []Job    // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
}

// setCore sets the core field of each job in page p to c.
func (p *JobsPage) setCore(c *Core) {
	for i := range p.Jobs {
		p.Jobs[i].setCore(c)
	}
}

// RequiredJobsPage retrieves a page of jobs required by job j.
func (j Job) RequiredJobsPage(ctx context.Context, pa PageArgs) (JobsPage, error) {
	p, err := j.c.p.GetJobsByID(ctx, pa, j.WorkflowID, j.Requires)
	p.setCore(j.c)
	return p, err
}
