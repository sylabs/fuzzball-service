// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"time"

	"github.com/blang/semver"
	"github.com/sylabs/fuzzball-service/internal/pkg/token"
)

// BuildInfo contains build information about the core.
type BuildInfo struct {
	GitVersion   *semver.Version
	GitCommit    *string
	GitTreeState *string
	BuiltAt      *time.Time
	GoVersion    string
	Compiler     string
	Platform     string
}

// GetBuildInfo returns build information about the core.
func (c *Core) GetBuildInfo(ctx context.Context) (BuildInfo, error) {
	if _, ok := token.FromContext(ctx); !ok {
		return BuildInfo{}, ErrNotAuthenticated
	}
	return c.bi, nil
}
