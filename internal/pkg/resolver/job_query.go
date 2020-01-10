// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

// Job returns a job resolver.
func (r Resolver) Job(ctx context.Context, args struct {
	ID string
}) (*JobResolver, error) {
	j, err := r.p.GetJob(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &JobResolver{&j}, nil
}
