// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

// SystemInfo returns pointer to SystemInfo resolver
func (r Resolver) SystemInfo(ctx context.Context) (*SystemInfoResolver, error) {
	return &SystemInfoResolver{r.si}, nil
}
