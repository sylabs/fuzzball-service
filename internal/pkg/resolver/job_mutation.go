// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

// CreateJob creates a new job.
func (r Resolver) CreateJob(ctx context.Context, args struct {
	Name string
}) (*JobResolver, error) {
	j, err := r.p.CreateJob(args.Name)
	if err != nil {
		return nil, err
	}
	return &JobResolver{j}, nil
}
