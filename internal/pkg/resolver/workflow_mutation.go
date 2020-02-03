// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/graph"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// CreateWorkflow creates a new workflow.
func (r Resolver) CreateWorkflow(ctx context.Context, args struct {
	Spec workflowSpec
}) (*WorkflowResolver, error) {
	// Create workflow document
	w, err := r.p.CreateWorkflow(ctx, model.Workflow{Name: args.Spec.Name})
	if err != nil {
		return nil, err
	}

	// iterate through jobSpecs and add them to the graph and a map by name for later
	g := graph.New()
	jobNameMapping := make(map[string]int)
	for i, js := range args.Spec.Jobs {
		requires := make([]string, 0)
		if js.Requires != nil {
			requires = *js.Requires
		}
		if err := g.AddVertex(js.Name, requires); err != nil {
			return nil, err
		}

		jobNameMapping[js.Name] = i
	}

	// ensure jobs are correctly referencing eachother semantically
	if err := g.Validate(); err != nil {
		return nil, err
	}

	// sort jobs by dependencies so we can insert them in
	// an order that allows for the parent IDs to have already been generated
	s, err := g.TopoSort()
	if err != nil {
		return nil, err
	}

	// create jobs in persistent storage
	var jobs []model.Job
	jobNameToID := make(map[string]string)
	for _, name := range s {
		// lookup job by name
		js := args.Spec.Jobs[jobNameMapping[name]]

		requires := []string{}
		if js.Requires != nil {
			// convert requires job name to job IDs
			for _, name := range *js.Requires {
				id, ok := jobNameToID[name]
				if !ok {
					return nil, fmt.Errorf("jobs created in invalid order")
				}

				requires = append(requires, id)
			}
		}

		j, err := r.p.CreateJob(ctx, model.Job{
			WorkflowID: w.ID,
			Name:       js.Name,
			Image:      js.Image,
			Command:    js.Command,
			Requires:   requires,
		})
		if err != nil {
			return nil, err
		}

		jobNameToID[j.Name] = j.ID
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
