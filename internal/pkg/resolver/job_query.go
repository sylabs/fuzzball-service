// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/mongodb"
)

type JobsArgs struct {
	ID   *string
	Name *string
}

func (r Resolver) Jobs(ctx context.Context, args JobsArgs) (*[]*JobResolver, error) {
	jobs, err := r.p.GetJobs(ctx, mongodb.JobsQueryArgs{
		ID:   args.ID,
		Name: args.Name,
	})
	if err != nil {
		return nil, err
	}

	var resolvers = make([]*JobResolver, 0, len(jobs))
	for n := 0; n<len(jobs); n++ {
		resolvers = append(resolvers, &JobResolver{&jobs[n]})
	}
	return &resolvers, nil
}
