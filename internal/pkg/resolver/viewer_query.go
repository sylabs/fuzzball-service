// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/core"
)

// Viewer returns the currently authenticated user.
func (r Resolver) Viewer(ctx context.Context) (*UserResolver, error) {
	return &UserResolver{
		u: &core.User{
			ID:    "507f1f77bcf86cd799439011",
			Login: "jimbob",
		},
		s: r.s,
	}, nil
}
