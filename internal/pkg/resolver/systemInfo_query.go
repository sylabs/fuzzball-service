// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"os"
)

// SystemInfo returns pointer to SystemInfo resolver
func (r Resolver) SystemInfo(ctx context.Context) (*SystemInfoResolver, error) {
	hostName, _ := os.Hostname()

	si := SystemInfo{
		HostName: hostName,
	}

	return &SystemInfoResolver{&si}, nil
}
