// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

// Schema validates the GraphQL schema and generates a Go file containing the schema.
func Schema() error {
	return generateSchema()
}

// ldFlags returns standard linker flags to pass to various Go commands.
func ldFlags() string {
	vals := []string{
		fmt.Sprintf("-X main.builtAt=%v", time.Now().UTC().Format(time.RFC3339)),
	}

	// Attempt to get git details.
	d, err := describeHead()
	if err == nil {
		vals = append(vals, fmt.Sprintf("-X main.gitCommit=%v", d.ref.Hash().String()))

		if d.isClean {
			vals = append(vals, "-X main.gitTreeState=clean")
		} else {
			vals = append(vals, "-X main.gitTreeState=dirty")
		}

		if v, err := getVersion(d); err != nil {
			logrus.WithError(err).Warn("failed to get version from git description")
		} else {
			vals = append(vals, fmt.Sprintf("-X main.gitVersion=%v", v.String()))
		}
	}

	return strings.Join(vals, " ")
}

// Build builds Fuzzball assets using `go build`.
func Build() error {
	mg.Deps(Schema)
	return sh.RunV(mg.GoCmd(), "build", "-ldflags", ldFlags(), "./...")
}

// Install installs Fuzzball assets using `go install`.
func Install() error {
	mg.Deps(Schema)
	return sh.RunV(mg.GoCmd(), "install", "-ldflags", ldFlags(), "./...")
}

// Run runs the Fuzzball server using `go run`.
func Run() error {
	mg.Deps(Schema)
	return sh.RunV(mg.GoCmd(), "run", "-ldflags", ldFlags(), "./cmd/server/")
}

// Test runs unit and integration tests using `go test`.
func Test() error {
	mg.Deps(Schema)
	return sh.RunV(mg.GoCmd(), "test", "-ldflags", ldFlags(), "-cover", "-race", "-tags=integration", "./...")
}

// UnitTest runs unit tests using `go test`.
func UnitTest() error {
	mg.Deps(Schema)
	return sh.RunV(mg.GoCmd(), "test", "-ldflags", ldFlags(), "-cover", "-race", "./...")
}
