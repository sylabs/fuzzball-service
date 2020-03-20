// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/fuzzball-service/internal/pkg/core"
)

// BuildInfoServicer is the interface by which build info is retrieved.
type BuildInfoServicer interface {
	GetBuildInfo(ctx context.Context) (core.BuildInfo, error)
}

// BuildInfoResolver resolves build info.
type BuildInfoResolver struct {
	bi core.BuildInfo
}

// GitVersion returns Semantic version info (if available).
func (r BuildInfoResolver) GitVersion() *SemanticVersionResolver {
	if v := r.bi.GitVersion; v != nil {
		return &SemanticVersionResolver{*v}
	}
	return nil
}

// GitCommit returns the specific git commit the component was built from (if available).
func (r BuildInfoResolver) GitCommit() *string {
	if c := r.bi.GitCommit; c != nil {
		return c
	}
	return nil
}

// GitTreeState returns the clean/dirty state of the git tree the component was built from (if
// available).
func (r BuildInfoResolver) GitTreeState() *string {
	if s := r.bi.GitTreeState; s != nil {
		return s
	}
	return nil
}

// BuiltAt returns the time at which the component was built (if available).
func (r BuildInfoResolver) BuiltAt() *graphql.Time {
	if t := r.bi.BuiltAt; t != nil {
		return &graphql.Time{Time: *t}
	}
	return nil
}

// GoVersion returns the version of Go the component utilizes.
func (r BuildInfoResolver) GoVersion() string {
	return r.bi.GoVersion
}

// Compiler returns the name of the compiler toolchain that built the component.
func (r BuildInfoResolver) Compiler() string {
	return r.bi.Compiler
}

// Platform returns the operating system and architecture of the component.
func (r BuildInfoResolver) Platform() string {
	return r.bi.Platform
}
