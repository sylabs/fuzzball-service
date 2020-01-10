// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	GetJob(ID string) (*model.Job, error)
}

// Job returns a job resolver.
func (r Resolver) Job(ctx context.Context, args struct {
	ID string
}) (*JobResolver, error) {
	j, err := r.p.GetJob(args.ID)
	if err != nil {
		return nil, err
	}
	return &JobResolver{j}, nil
}
