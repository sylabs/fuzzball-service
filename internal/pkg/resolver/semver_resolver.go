// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"strings"

	"github.com/blang/semver"
)

// SemanticVersionResolver resolves version information
type SemanticVersionResolver struct {
	v semver.Version
}

// Major returns the major version.
func (r SemanticVersionResolver) Major() int32 {
	return int32(r.v.Major)
}

// Minor returns the minor version.
func (r SemanticVersionResolver) Minor() int32 {
	return int32(r.v.Minor)
}

// Patch returns the patch version.
func (r SemanticVersionResolver) Patch() int32 {
	return int32(r.v.Patch)
}

// PreRelease returns pre-release info.
func (r SemanticVersionResolver) PreRelease() *string {
	var pre []string
	for _, p := range r.v.Pre {
		pre = append(pre, p.String())
	}
	if s := strings.Join(pre, "."); s != "" {
		return &s
	}
	return nil
}

// BuildMetadata returns build metadata.
func (r SemanticVersionResolver) BuildMetadata() *string {
	if s := strings.Join(r.v.Build, "."); s != "" {
		return &s
	}
	return nil
}
