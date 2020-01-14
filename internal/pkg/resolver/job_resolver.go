// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
	"github.com/sylabs/compute-service/internal/pkg/mongodb"
)

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	CreateJob(context.Context, model.Job) (model.Job, error)
	DeleteJob(context.Context, string) (model.Job, error)
	GetJobs(context.Context, mongodb.JobsQueryArgs) ([]model.Job, error)
}

// JobResolver resolves a job.
type JobResolver struct {
	j *model.Job
}

// ID resolves the job ID.
func (r *JobResolver) ID() graphql.ID {
	return graphql.ID(r.j.ID)
}

// Name resolves the job name.
func (r *JobResolver) Name() string {
	return r.j.Name
}
