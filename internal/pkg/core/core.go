// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"errors"

	"github.com/sylabs/fuzzball-service/internal/pkg/token"
)

// Persister is the interface by which all data is persisted.
type Persister interface {
	WorkflowPersister
	JobPersister
	VolumePersister
}

// IOFetcher is the interface where IO data is retrieved.
type IOFetcher interface {
	JobOutputFetcher
}

// Scheduler is the interface by which all workflows are scheduled.
type Scheduler interface {
	AddWorkflow(context.Context, Workflow, []Job, map[string]Volume) error
}

// Core represents core business logic.
type Core struct {
	p Persister
	f IOFetcher
	s Scheduler
}

// New creates a new core.
func New(p Persister, f IOFetcher, s Scheduler) (*Core, error) {
	return &Core{p: p, f: f, s: s}, nil
}

// Viewer returns the user associated with ctx.
func (c *Core) Viewer(ctx context.Context) (User, error) {
	t, ok := token.FromContext(ctx)
	if !ok {
		return User{}, errors.New("viewer not logged in")
	}
	tc := t.Claims()

	u := User{
		ID:    tc.UserID,
		Login: tc.Subject,
	}
	u.setCore(c)
	return u, nil
}

// WorkflowSpec represents a workflow specification.
type WorkflowSpec struct {
	Name    string        `bson:"name"`
	Jobs    []jobSpec     `bson:"jobs"`
	Volumes *[]volumeSpec `bson:"volumes"`
}

type jobSpec struct {
	Name     string                   `bson:"name"`
	Image    string                   `bson:"image"`
	Command  []string                 `bson:"command"`
	Requires *[]string                `bson:"requires"`
	Volumes  *[]volumeRequirementSpec `bson:"volumes"`
}

type volumeRequirementSpec struct {
	Name     string
	Location string
}

// CreateWorkflow creates a new workflow. If an ID is provided in w, it is ignored and replaced
// with a unique identifier in the returned workflow.
func (c *Core) CreateWorkflow(ctx context.Context, s WorkflowSpec) (Workflow, error) {
	// TODO: authorization
	w, err := c.p.CreateWorkflow(ctx, Workflow{Name: s.Name})
	if err != nil {
		return Workflow{}, err
	}

	volumes, err := createVolumes(ctx, c.p, w, s.Volumes)
	if err != nil {
		return Workflow{}, err
	}

	// Jobs must be created after volumes to allow them to reference
	// generated volume IDs
	jobs, err := c.createJobs(ctx, w, volumes, s.Jobs)
	if err != nil {
		return Workflow{}, err
	}

	// Schedule the workflow.
	if err := c.s.AddWorkflow(ctx, w, jobs, volumes); err != nil {
		return Workflow{}, err
	}

	w.setCore(c)
	return w, err
}

// DeleteWorkflow deletes a workflow by ID. If the supplied ID is not valid, or there there is not
// a workflow with a matching ID in the database, an error is returned.
func (c *Core) DeleteWorkflow(ctx context.Context, id string) (Workflow, error) {
	// TODO: authorization
	w, err := c.p.DeleteWorkflow(ctx, id)
	if err != nil {
		return Workflow{}, err
	}

	err = c.p.DeleteJobsByWorkflowID(ctx, w.ID)
	if err != nil {
		return Workflow{}, err
	}

	err = c.p.DeleteVolumesByWorkflowID(ctx, w.ID)
	if err != nil {
		return Workflow{}, err
	}

	w.setCore(c)
	return w, nil
}

// GetWorkflow retrieves a workflow by ID. If the supplied ID is not valid, or there there is not a
// workflow with a matching ID in the database, an error is returned.
func (c *Core) GetWorkflow(ctx context.Context, id string) (Workflow, error) {
	// TODO: authorization
	w, err := c.p.GetWorkflow(ctx, id)
	w.setCore(c)
	return w, err
}
