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

func (r Resolver) Jobs(ctx context.Context, args struct {
	ID   *string
	Name *string
}) (*[]*JobResolver, error) {
	filterSpec := map[string]string{}

	if args.ID != nil {
		filterSpec["ID"] = *args.ID
	}

	if args.Name != nil {
		filterSpec["name"] = *args.Name
	}

	jobs, err := r.p.GetJobs(ctx, filterSpec)
	if err != nil {
		return nil, err
	}

	var resolvers = make([]*JobResolver, 0, len(jobs))
	for n := 0; n < len(jobs); n++ {
		resolvers = append(resolvers, &JobResolver{&jobs[n]})
	}
	return &resolvers, nil
}
