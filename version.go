// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"errors"
	"strings"

	"github.com/blang/semver"
)

// getVersion returns a semantic version based on d.
func getVersion(d *gitDescription) (semver.Version, error) {
	if d.tag == nil {
		return semver.Version{}, errors.New("no semver tags found")
	}

	v, err := semver.Parse(strings.TrimPrefix(d.tag.Name, "v"))
	if err != nil {
		return semver.Version{}, err
	}

	// If this version wasn't tagged directly, modify tag.
	if d.n > 0 {
		if len(v.Pre) == 0 {
			// The tag is not a pre-release version. Bump the patch version and add a pre-release
			// of alpha.0. Semantically, this indicates this is pre-alpha.1, which would normally
			// be the first alpha version.
			v.Patch += 1
			v.Pre = append(v.Pre, semver.PRVersion{VersionStr: "alpha"})
			v.Pre = append(v.Pre, semver.PRVersion{VersionNum: 0, IsNum: true})
		}

		// Append devel.N to pre-release version. For example, if the tag is 0.1.2-alpha.1, tag as
		// 0.1.2-alpha.1.devel.3. Semantically, this indicates this version is between alpha.1 and
		// alpha.2.
		v.Pre = append(v.Pre, semver.PRVersion{VersionStr: "devel"})
		v.Pre = append(v.Pre, semver.PRVersion{VersionNum: d.n, IsNum: true})
	}

	return v, nil
}
