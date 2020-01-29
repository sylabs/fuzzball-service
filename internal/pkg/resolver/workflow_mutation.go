// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// CreateWorkflow creates a new workflow.
func (r Resolver) CreateWorkflow(ctx context.Context, args struct {
	Spec model.WorkflowSpec
}) (*WorkflowResolver, error) {
	// Create workflow document
	w, err := r.p.CreateWorkflow(ctx, model.Workflow{Name: args.Spec.Name})
	if err != nil {
		return nil, err
	}

	var jobs []model.Job
	for _, js := range args.Spec.Jobs {
		j, err := r.p.CreateJob(ctx, model.Job{
			WorkflowID: w.ID,
			Name:       js.Name,
			Image:      js.Image,
			Command:    js.Command,
		})
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	// Schedule the workflow.
	if err := r.s.AddWorkflow(ctx, w, jobs); err != nil {
		return nil, err
	}

	return &WorkflowResolver{w, r.p}, nil
}

// DeleteWorkflow deletes a workflow.
func (r Resolver) DeleteWorkflow(ctx context.Context, args struct {
	ID string
}) (*WorkflowResolver, error) {
	w, err := r.p.DeleteWorkflow(ctx, args.ID)
	if err != nil {
		return nil, err
	}

	err = r.p.DeleteJobsByWorkflowID(ctx, w.ID)
	if err != nil {
		return nil, err
	}

	return &WorkflowResolver{w, r.p}, nil
}
