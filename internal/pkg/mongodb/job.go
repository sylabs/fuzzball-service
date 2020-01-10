// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// CreateJob creates a new job.
func (c *Connection) CreateJob(name string) (*model.Job, error) {
	// TODO
	return &model.Job{
		ID:   "1234",
		Name: name,
	}, nil
}

// GetJob retrieves a job by ID.
func (c *Connection) GetJob(id string) (*model.Job, error) {
	// TODO
	return &model.Job{
		ID:   "1234",
		Name: "jobName",
	}, nil
}
