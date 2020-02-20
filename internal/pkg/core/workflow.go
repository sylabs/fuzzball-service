// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/graph"
)

// WorkflowPersister is the interface by which workflows are persisted.
type WorkflowPersister interface {
	CreateWorkflow(context.Context, Workflow) (Workflow, error)
	DeleteWorkflow(context.Context, string) (Workflow, error)
	GetWorkflow(context.Context, string) (Workflow, error)
	GetWorkflows(context.Context, PageArgs) (WorkflowsPage, error)
}

// Workflow represents a workflow.
type Workflow struct {
	ID     string `bson:"_id,omitempty"`
	Name   string `bson:"name"`
	Status string `bson:"status"`
}

// WorkflowsPage represents a page of workflows resulting from a query, and associated metadata.
type WorkflowsPage struct {
	Workflows  []Workflow // Slice of results.
	PageInfo   PageInfo   // Information to aid in pagination.
	TotalCount int        // Identifies the total count of items in the connection.
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

	return w, err
}

// DeleteWorkflow deletes a workflow by ID. If the supplied ID is not valid, or there there is not
// a workflow with a matching ID in the database, an error is returned.
func (c *Core) DeleteWorkflow(ctx context.Context, id string) (Workflow, error) {
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

	return w, nil
}

// GetWorkflow retrieves a workflow by ID. If the supplied ID is not valid, or there there is not a
// workflow with a matching ID in the database, an error is returned.
func (c *Core) GetWorkflow(ctx context.Context, id string) (w Workflow, err error) {
	return c.p.GetWorkflow(ctx, id)
}

// GetWorkflows returns a list of all workflows.
func (c *Core) GetWorkflows(ctx context.Context, pa PageArgs) (p WorkflowsPage, err error) {
	return c.p.GetWorkflows(ctx, pa)
}

func (c *Core) createJobs(ctx context.Context, w Workflow, volumes map[string]Volume, specs []jobSpec) ([]Job, error) {
	// iterate through jobSpecs and add them to the graph and a map by name for later
	g := graph.New()
	jobNameMapping := make(map[string]int)
	for i, js := range specs {
		// check job spec for invalid volume references
		if js.Volumes != nil {
			for _, v := range *js.Volumes {
				if _, ok := volumes[v.Name]; !ok {
					return nil, fmt.Errorf("job %q references nonexistant volume %q", js.Name, v.Name)
				}
			}
		}

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
	var jobs []Job
	jobNameToID := make(map[string]string)
	for _, name := range s {
		// lookup job by name
		js := specs[jobNameMapping[name]]

		// construct list of required job IDs
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

		// construct list of required volume IDs
		volumeReqs := []VolumeRequirement{}
		if js.Volumes != nil {
			for _, v := range *js.Volumes {
				volumeReqs = append(volumeReqs, VolumeRequirement{
					Name:     v.Name,
					Location: v.Location,
					VolumeID: volumes[v.Name].ID,
				})
			}
		}

		j, err := c.CreateJob(ctx, Job{
			WorkflowID: w.ID,
			Name:       js.Name,
			Image:      js.Image,
			Command:    js.Command,
			Requires:   requires,
			Volumes:    volumeReqs,
		})
		if err != nil {
			return nil, err
		}

		jobNameToID[j.Name] = j.ID
		jobs = append(jobs, j)
	}

	return jobs, nil
}
