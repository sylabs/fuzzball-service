// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// GetJob retrieves a job by ID.
func (c *Connection) GetJob(ID string) (*model.Job, error) {
	// TODO
	return &model.Job{
		ID:   "1234",
		Name: "jobName",
	}, nil
}
