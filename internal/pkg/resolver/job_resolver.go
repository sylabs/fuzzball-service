// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	CreateJob(name string) (*model.Job, error)
	GetJob(id string) (*model.Job, error)
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
