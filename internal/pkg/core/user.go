// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
)

// User represents a user.
type User struct {
	ID    string `bson:"_id,omitempty"` // Unique user ID.
	Login string `bson:"login"`         // The username used to login.

	c *Core // Used internally for lazy loading.
}

// setCore sets the core of w to c.
func (u *User) setCore(c *Core) {
	u.c = c
}

// WorkflowsPage retrieves a page of workflows created by user u.
func (u User) WorkflowsPage(ctx context.Context, pa PageArgs) (WorkflowsPage, error) {
	p, err := u.c.p.GetWorkflows(ctx, pa)
	p.setCore(u.c)
	return p, err
}

// JobsPage retrieves a page of jobs created by user u.
func (u User) JobsPage(ctx context.Context, pa PageArgs) (JobsPage, error) {
	p, err := u.c.p.GetJobs(ctx, pa)
	p.setCore(u.c)
	return p, err
}

// VolumesPage retrieves a page of volumes created by user u.
func (u User) VolumesPage(ctx context.Context, pa PageArgs) (VolumesPage, error) {
	p, err := u.c.p.GetVolumes(ctx, pa)
	p.setCore(u.c)
	return p, err
}
