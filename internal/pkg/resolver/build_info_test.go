// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/blang/semver"
	"github.com/sylabs/fuzzball-service/internal/pkg/core"
	"github.com/sylabs/fuzzball-service/internal/pkg/schema"
)

func TestBuildInfo(t *testing.T) {
	gitTreeState := "clean"
	gitCommit := "50b3625811b4a5a6ecd3fcdf3a0180cf2908b651"
	builtAt := time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC)

	tests := []struct {
		name         string
		gitVersion   *semver.Version
		gitCommit    *string
		gitTreeState *string
		builtAt      *time.Time
	}{
		{
			name: "Version0.1.2",
			gitVersion: &semver.Version{
				Minor: 1,
				Patch: 2,
			},
		},
		{
			name: "Version0.1.2-alpha.1",
			gitVersion: &semver.Version{
				Minor: 1,
				Patch: 2,
				Pre: []semver.PRVersion{
					{VersionStr: "alpha"},
					{VersionNum: 1, IsNum: true},
				},
			},
		},
		{
			name: "Version1.2.3+four.five",
			gitVersion: &semver.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: []string{
					"four",
					"five",
				},
			},
		},
		{
			name:         "GitTreeState",
			gitTreeState: &gitTreeState,
		},
		{
			name:      "GitCommit",
			gitCommit: &gitCommit,
		},
		{
			name:    "BuiltAt",
			builtAt: &builtAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := [](func(*core.Core) error){
				core.OptGoVersionOverride("go1.14"),
				core.OptCompilerOverride("gc"),
				core.OptPlatformOverride("linux/amd64"),
			}
			if v := tt.gitVersion; v != nil {
				opts = append(opts, core.OptGitVersion(*v))
			}
			if c := tt.gitCommit; c != nil {
				opts = append(opts, core.OptGitCommit(*c))
			}
			if s := tt.gitTreeState; s != nil {
				opts = append(opts, core.OptGitTreeState(*s))
			}
			if t := tt.builtAt; t != nil {
				opts = append(opts, core.OptBuiltAt(*t))
			}

			mc, err := getMockCore(mockCore{}, opts...)
			if err != nil {
				t.Fatal(err)
			}

			s, err := schema.Get(&Resolver{s: mc})
			if err != nil {
				t.Fatal(err)
			}

			q := `
			query OpName()
			{
			  serverBuildInfo {
			    gitVersion {
			      major
			      minor
			      patch
			      preRelease
			      buildMetadata
				}
			    gitCommit
			    gitTreeState
			    builtAt
			    goVersion
			    compiler
			    platform
			  }
			}`

			res := s.Exec(context.Background(), q, "", nil)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
