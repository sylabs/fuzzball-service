// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

// Workflow returns a workflow resolver.
func (r Resolver) Workflow(ctx context.Context, args struct {
	ID string
}) (*WorkflowResolver, error) {
	j, err := r.p.GetWorkflow(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	return &WorkflowResolver{j, r.p}, nil
}
