// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"fmt"
	"path/filepath"

	"github.com/go-bindata/go-bindata/v3"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// ldFlags returns standard linker flags to pass to various Go commands.
func ldFlags() string {
	return fmt.Sprintf("-X main.version=%s", version())
}

// generateSchema generates a Go file containing the GraphQL schema.
func generateSchema() error {
	c := bindata.NewConfig()
	c.NoMetadata = true
	c.Package = "schema"
	c.Input = []bindata.InputConfig{{Path: filepath.Join("api", "graphql-schema")}}
	c.Output = filepath.Join("internal", "pkg", "schema", "bindata.go")

	return bindata.Translate(c)
}

// Build builds Fuzzball assets using `go build`.
func Build() error {
	mg.Deps(generateSchema)
	return sh.RunV(mg.GoCmd(), "build", "-ldflags", ldFlags(), "./...")
}

// Install installs Fuzzball assets using `go install`.
func Install() error {
	mg.Deps(generateSchema)
	return sh.RunV(mg.GoCmd(), "install", "-ldflags", ldFlags(), "./...")
}

// Run runs the Fuzzball server using `go run`.
func Run() error {
	mg.Deps(generateSchema)
	return sh.RunV(mg.GoCmd(), "run", "-ldflags", ldFlags(), "./cmd/server/")
}

// Test runs unit and integration tests using `go test`.
func Test() error {
	mg.Deps(generateSchema)
	return sh.RunV(mg.GoCmd(), "test", "-ldflags", ldFlags(), "-cover", "-race", "-tags=integration", "./...")
}

// UnitTest runs unit tests using `go test`.
func UnitTest() error {
	mg.Deps(generateSchema)
	return sh.RunV(mg.GoCmd(), "test", "-ldflags", ldFlags(), "-cover", "-race", "./...")
}
