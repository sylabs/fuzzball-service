// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

// Viewer returns the currently authenticated user.
func (r Resolver) Viewer(ctx context.Context) (*UserResolver, error) {
	u, err := r.s.Viewer(ctx)
	if err != nil {
		return nil, err
	}
	return &UserResolver{&u}, nil
}
