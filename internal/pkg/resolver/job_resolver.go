// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/core"
)

// JobResolver resolves a job.
type JobResolver struct {
	j core.Job
}

// ID resolves the job ID.
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
func (r *JobResolver) CreatedBy(ctx context.Context) (*UserResolver, error) {
	u, err := r.j.CreatedBy(ctx)
	if err != nil {
		return nil, err
	}
	return &UserResolver{u: &u}, nil
}

// CreatedAt resolves when the job was created.
func (r *JobResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.j.CreatedAt}
}

// StartedAt returns when the job started, if it has started.
func (r *JobResolver) StartedAt() *graphql.Time {
	if t := r.j.StartedAt; t != nil {
		return &graphql.Time{Time: *t}
	}
	return nil
}

// FinishedAt returns when the job finished, if it has finished.
func (r *JobResolver) FinishedAt() *graphql.Time {
	if t := r.j.FinishedAt; t != nil {
		return &graphql.Time{Time: *t}
	}
	return nil
}

// Status resolves the state of the job.
func (r *JobResolver) Status() string {
	return r.j.Status
}

// ExitCode resolves the exit status process that ran the job.
func (r *JobResolver) ExitCode() *int32 {
	if r.j.ExitCode != nil {
		i := int32(*r.j.ExitCode)
		return &i
	}
	return nil
}

// Output resolves the captured Stdout/Stderr of the job.
func (r *JobResolver) Output() (string, error) {
	return r.j.GetOutput()
}

// Requires looks up jobs that need to be executed before the current one.
func (r *JobResolver) Requires(ctx context.Context, args pageArgs) (*JobConnectionResolver, error) {
	p, err := r.j.RequiredJobsPage(ctx, convertPageArgs(args))
	if err != nil {
		return nil, err
	}
	return &JobConnectionResolver{p}, nil
}
