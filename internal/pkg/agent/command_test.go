// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build linux darwin

package agent

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	ctx := context.Background()
	expiredCtx, cancel := context.WithDeadline(ctx, time.Now().Add(-time.Hour))
	defer cancel()

	catPath, err := exec.LookPath("cat")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		ctx   context.Context
		path  string
		args  []string
		env   []string
		dir   string
		stdin io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{"CatEmpty", args{
			ctx:  ctx,
			path: catPath,
		}, "", "", false},
		{"CatStdin", args{
			ctx:   ctx,
			path:  catPath,
			stdin: strings.NewReader("hello"),
		}, "hello", "", false},
		{"CatFile", args{
			ctx:  ctx,
			path: catPath,
			args: []string{path.Join("testdata", "hello.txt")},
		}, "hello", "", false},
		{"ExpiredContext", args{
			ctx:  expiredCtx,
			path: catPath,
		}, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			_, err := runCommand(tt.args.ctx, tt.args.path, tt.args.args, tt.args.env, tt.args.dir, tt.args.stdin, stdout, stderr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error %v, wantErr %v", err, tt.wantErr)
			}
			if gotStdout := stdout.String(); gotStdout != tt.wantStdout {
				t.Errorf("got stdout %v, want %v", gotStdout, tt.wantStdout)
			}
			if gotStderr := stderr.String(); gotStderr != tt.wantStderr {
				t.Errorf("got stderr %v, want %v", gotStderr, tt.wantStderr)
			}
		})
	}
}
