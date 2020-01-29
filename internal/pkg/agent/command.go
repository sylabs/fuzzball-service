// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package agent

import (
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

// runCommand runs the command specified by name, with arguments args, with stdin, stdout and
// stderr connected as one would expect.
func runCommand(ctx context.Context, path string, args, env []string, dir string, stdin io.Reader, stdout, stderr io.Writer) (*os.ProcessState, error) {
	cmd := exec.CommandContext(ctx, path, args...)
	cmd.Env = env
	cmd.Dir = dir
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	// Start the process.
	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Log process start and finish.
	log := logrus.WithFields(logrus.Fields{
		"path": path,
		"args": args,
		"env":  env,
		"dir":  dir,
		"pid":  cmd.Process.Pid,
	})
	log.Print("command started")
	defer func(t time.Time, cmd *exec.Cmd) {
		log.WithFields(logrus.Fields{
			"wallTime":   time.Since(t),
			"userTime":   cmd.ProcessState.UserTime(),
			"systemTime": cmd.ProcessState.SystemTime(),
			"exitCode":   cmd.ProcessState.ExitCode(),
		}).Print("command finished")
	}(startTime, cmd)

	// Wait for process to finish.
	err := cmd.Wait()
	return cmd.ProcessState, err
}
