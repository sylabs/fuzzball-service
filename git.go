// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

// +build mage

package main

import (
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
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

type gitDescription struct {
	isClean bool                // if true, the git working tree has local modifications
	ref     *plumbing.Reference // reference being described
	tag     *object.Tag         // nearest semver tag reachable from ref (or nil if none found)
	n       uint64              // number of commits between nearest semver tag and ref (if tag is non-nil)
}

// describe returns a gitDescription of ref.
func describe(r *git.Repository, ref *plumbing.Reference) (*gitDescription, error) {
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := w.Status()
	if err != nil {
		return nil, err
	}

	// Get version tags.
	tags, err := getVersionTags(r)
	if err != nil {
		return nil, err
	}

	// Get commit log.
	logIter, err := r.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		From:  ref.Hash(),
	})
	if err != nil {
		return nil, err
	}

	// Iterate through commit log until we find a matching tag.
	var tag *object.Tag
	var n uint64
	err = logIter.ForEach(func(c *object.Commit) error {
		if t, ok := tags[c.Hash]; ok {
			tag = t
		}
		if tag != nil {
			return storer.ErrStop
		}
		n++
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &gitDescription{
		isClean: status.IsClean(),
		ref:     ref,
		tag:     tag,
		n:       n,
	}, nil
}

// describeHead returns a gitDescription of HEAD.
func describeHead() (*gitDescription, error) {
	// Open git repo.
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, err
	}

	// Get HEAD commit.
	head, err := r.Head()
	if err != nil {
		return nil, err
	}

	return describe(r, head)
}
