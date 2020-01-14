// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"runtime"
	"strconv"
)

// NewCapabilities returns a capability resolver.
func NewCapabilities(ctx context.Context) (*[]*CapabilityResolver, error) {
	var resolvers []*CapabilityResolver

	c1 := CapabilityResolver{
		&Capability{
			key:   "CPUArchitecture",
			value: runtime.GOARCH,
		},
	}

	resolvers = append(resolvers, &c1)

	c2 := CapabilityResolver{
		&Capability{
			key:   "NumCPU",
			value: strconv.Itoa(runtime.NumCPU()),
		},
	}

	resolvers = append(resolvers, &c2)

	c3 := CapabilityResolver{
		&Capability{
			key:   "GPU",
			value: "false",
		},
	}

	resolvers = append(resolvers, &c3)

	return &resolvers, nil
}
