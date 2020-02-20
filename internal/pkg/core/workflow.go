// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// WorkflowPersister is the interface by which workflows are persisted.
type WorkflowPersister interface {
	CreateWorkflow(context.Context, model.Workflow) (model.Workflow, error)
	DeleteWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflow(context.Context, string) (model.Workflow, error)
	GetWorkflows(context.Context, model.PageArgs) (model.WorkflowsPage, error)
}

// WorkflowSpec represents a workflow specification.
type WorkflowSpec struct {
	Name    string        `bson:"name"`
	Jobs    []jobSpec     `bson:"jobs"`
	Volumes *[]volumeSpec `bson:"volumes"`
}

// CreateWorkflow creates a new workflow. If an ID is provided in w, it is ignored and replaced
// with a unique identifier in the returned workflow.
func (c *Core) CreateWorkflow(ctx context.Context, s WorkflowSpec) (model.Workflow, error) {
	w, err := c.p.CreateWorkflow(ctx, model.Workflow{Name: s.Name})
	if err != nil {
		return model.Workflow{}, err
	}

	volumes, err := createVolumes(ctx, c.p, w, s.Volumes)
	if err != nil {
		return model.Workflow{}, err
	}

	// Jobs must be created after volumes to allow them to reference
	// generated volume IDs
	jobs, err := createJobs(ctx, c.p, w, volumes, s.Jobs)
	if err != nil {
		return model.Workflow{}, err
	}

	// Schedule the workflow.
	if err := c.s.AddWorkflow(ctx, w, jobs, volumes); err != nil {
		return model.Workflow{}, err
	}

	return w, err
}

// DeleteWorkflow deletes a workflow by ID. If the supplied ID is not valid, or there there is not
// a workflow with a matching ID in the database, an error is returned.
func (c *Core) DeleteWorkflow(ctx context.Context, id string) (model.Workflow, error) {
	w, err := c.p.DeleteWorkflow(ctx, id)
	if err != nil {
		return model.Workflow{}, err
	}

	err = c.p.DeleteJobsByWorkflowID(ctx, w.ID)
	if err != nil {
		return model.Workflow{}, err
	}

	err = c.p.DeleteVolumesByWorkflowID(ctx, w.ID)
	if err != nil {
		return model.Workflow{}, err
	}

	return w, nil
}

// GetWorkflow retrieves a workflow by ID. If the supplied ID is not valid, or there there is not a
// workflow with a matching ID in the database, an error is returned.
func (c *Core) GetWorkflow(ctx context.Context, id string) (w model.Workflow, err error) {
	return c.p.GetWorkflow(ctx, id)
}

// GetWorkflows returns a list of all workflows.
func (c *Core) GetWorkflows(ctx context.Context, pa model.PageArgs) (p model.WorkflowsPage, err error) {
	return c.p.GetWorkflows(ctx, pa)
}
