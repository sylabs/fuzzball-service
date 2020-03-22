// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

// goldenPath returns the path of the golden file corresponding to name.
func goldenPath(name string) string {
	// Replace test name separator with OS-specific path separator.
	name = path.Join(strings.Split(name, "/")...)
	return path.Join("testdata", name) + ".golden"
}

// updateGolden writes b to a golden file associated with name.
func updateGolden(name string, b []byte) error {
	p := goldenPath(name)
	if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(p, b, 0644)
}

// verifyGolden compares b to the contents golden file associated with name.
func verifyGolden(name string, b []byte) error {
	if *update {
		if err := updateGolden(name, b); err != nil {
			return err
		}
	}
	g, err := ioutil.ReadFile(goldenPath(name))
	if err != nil {
		return err
	}

	if !bytes.Equal(b, g) {
		return errors.New("output does not match golden file")
	}
	return nil
}

func verifyGoldenJSON(name string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return verifyGolden(name, b)
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	mc, err := getMockCore(mockCore{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := New(mc, OAuth2Configuration{}); err != nil {
		t.Fatal(err)
	}
}
