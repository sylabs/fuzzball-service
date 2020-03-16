// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/sirupsen/logrus"
)

// getVersionTags returns a map of commit hashes to tags.
func getVersionTags(r *git.Repository) (map[plumbing.Hash]*object.Tag, error) {
	// Compile regex to select version tags.
	re, err := regexp.Compile(`^v[0-9]+(\.[0-9]+)?(\.[0-9]+)?` +
		`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)
	if err != nil {
		return nil, err
	}

	// Get a list of tags. Note that we cannot use r.TagObjects() directly, since that returns
	// objects that are not referenced (for example, deleted tags.)
	tagIter, err := r.Tags()
	if err != nil {
		return nil, err
	}

	// Iterate through tags, selecting tags that match regex.
	tags := make(map[plumbing.Hash]*object.Tag)
	err = tagIter.ForEach(func(ref *plumbing.Reference) error {
		if name := ref.Name().Short(); re.MatchString(name) {
			t, err := r.TagObject(ref.Hash())
			if err != nil {
				return err
			}
			tags[t.Target] = t
		}
		return nil
	})
	return tags, err
}

// describe returns the semver tag closest to the current commit, and the number of commits since
// the tag.
func describe(r *git.Repository, from plumbing.Hash) (*object.Tag, int, error) {
	// Get version tags.
	tags, err := getVersionTags(r)
	if err != nil {
		return nil, 0, err
	}

	// Get commit log.
	logIter, err := r.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		From:  from,
	})
	if err != nil {
		return nil, 0, err
	}

	// Iterate through commit log until we find a matching tag.
	var tag *object.Tag
	var count int
	err = logIter.ForEach(func(c *object.Commit) error {
		if t, ok := tags[c.Hash]; ok {
			tag = t
		}
		if tag != nil {
			return storer.ErrStop
		}
		count++
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	// If no tag found, return error. Otherwise return tag and distance.
	if tag == nil {
		return nil, 0, errors.New("descriptive version tag not found")
	}
	return tag, count, err
}

// version attempts to return a semantic version based on git tags. If the current commit is tagged
// with a semver tag (ex. v0.1.0), the version is returned (0.1.0). If the current commit is not
// tagged, the most recent tag is used, and the distance from the tag and the commit hash are
// appended to the version (ex. v0.1.0+14-g3b038b67).
func version() string {
	// Open git repo.
	r, err := git.PlainOpen(".")
	if err != nil {
		logrus.WithError(err).Warn("mage: failed to open git repo")
		return "unknown"
	}

	// Get HEAD commit.
	head, err := r.Head()
	if err != nil {
		logrus.WithError(err).Warn("mage: failed to get HEAD")
		return "unknown"
	}
	headHash := head.Hash().String()[:8]

	// Get descriptive tag and distance from the tag.
	tag, n, err := describe(r, head.Hash())
	if err != nil {
		logrus.WithError(err).Warn("mage: failed to describe git commit")
		return fmt.Sprintf("0.0.0+0-g%v", headHash)
	}

	// Build version.
	if n > 0 {
		return fmt.Sprintf("%v+%v-g%v", tag.Name, n, headHash)
	}
	return tag.Name
}
