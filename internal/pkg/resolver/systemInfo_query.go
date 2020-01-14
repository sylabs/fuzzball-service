// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"os"
	"runtime"

	"github.com/pbnjay/memory"
)

// SystemInfo returns pointer to SystemInfo resolver
func (r Resolver) SystemInfo(ctx context.Context) (*SystemInfoResolver, error) {
	hostName, _ := os.Hostname()

	si := SystemInfo{
		HostName:        hostName,
		CPUArchitecture: runtime.GOARCH,
		OSPlatform:      runtime.GOOS,
		Memory:          memory.TotalMemory() / 1024 / 1024, // megabytes
	}

	return &SystemInfoResolver{&si}, nil
}
