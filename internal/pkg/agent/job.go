// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"context"
	"os/exec"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// RunJob runs the specified job, returning the process exitCode.
func RunJob(ctx context.Context, j model.Job) (int, error) {
	// Locate Singularity in PATH.
	path, err := exec.LookPath("singularity")
	if err != nil {
		return 0, err
	}

	// Build up arguments to Singularity.
	args := []string{
		"exec",
		j.Image,
	}
	args = append(args, j.Command...)

	// Run Singularity.
	s, err := runCommand(ctx, path, args, []string{}, "", nil, nil, nil)
	if err != nil {
		if s != nil {
			return s.ExitCode(), err
		}
		return 0, err
	}
	return s.ExitCode(), nil
}
