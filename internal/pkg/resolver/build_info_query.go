// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import "context"

// ServerBuildInfo returns a server build info.
func (r Resolver) ServerBuildInfo(ctx context.Context) (BuildInfoResolver, error) {
	bi, err := r.s.GetBuildInfo(ctx)
	if err != nil {
		return BuildInfoResolver{}, err
	}
	return BuildInfoResolver{bi}, nil
}
