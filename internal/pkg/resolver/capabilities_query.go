// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

// NewCapabilities returns a capability resolver.
func NewCapabilities(ctx context.Context) (*[]*CapabilityResolver, error) {
	var resolvers []*CapabilityResolver

	c := CapabilityResolver{
		&Capability{
			key:   "GPU",
			value: "false",
		},
	}

	resolvers = append(resolvers, &c)

	return &resolvers, nil
}
