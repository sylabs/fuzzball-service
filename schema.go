// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/graph-gophers/graphql-go"
)

var (
	schemaInputDir     = filepath.Join("api", "graphql-schema")
	schemaTemplatePath = filepath.Join("internal", "pkg", "schema", "bindata.go.tmpl")
	schemaOutputPath   = filepath.Join("internal", "pkg", "schema", "bindata.go")
)

// readFiles walks the directory specified by root, writing the contents of each file to w.
func readFiles(root string, w io.Writer) error {
	return filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Add terminating newline, if not present.
		if nl := []byte("\n"); !bytes.HasSuffix(b, nl) {
			b = append(b, nl...)
		}

		_, err = w.Write(b)
		return err
	})
}

// writeCompressedSchema validates the schema, compresses it and writes it to w.
func writeCompressedSchema(w io.Writer) error {
	// Read the schema into buffer.
	b := &bytes.Buffer{}
	if err := readFiles(schemaInputDir, b); err != nil {
		return err
	}

	// Validate the schema.
	if _, err := graphql.ParseSchema(b.String(), nil, graphql.UseStringDescriptions()); err != nil {
		return err
	}

	// Write schema in compressed format.
	zw := zlib.NewWriter(w)
	defer zw.Close()

	// Write out schema.
	_, err := io.Copy(zw, b)
	return err
}

// generateSchema generates a Go file that exposes the schema.
func generateSchema() error {
	// Get schema.
	b := &bytes.Buffer{}
	if err := writeCompressedSchema(b); err != nil {
		return err
	}

	// Open file to write generated Go code into.
	f, err := os.OpenFile(schemaOutputPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write schema template.
	t, err := template.ParseFiles(schemaTemplatePath)
	if err != nil {
		return err
	}
	args := struct {
		Data        []byte
		GeneratedBy string
	}{
		b.Bytes(),
		"mage",
	}
	return t.Execute(f, args)
}
