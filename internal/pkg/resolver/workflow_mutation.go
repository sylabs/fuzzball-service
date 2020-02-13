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

	volumes, err := createVolumes(ctx, r.p, w, args.Spec.Volumes)
	if err != nil {
		return nil, err
	}

	// Jobs must be created after volumes to allow them to reference
	// generated volume IDs
	jobs, err := createJobs(ctx, r.p, w, volumes, args.Spec.Jobs)
	if err != nil {
		return nil, err
	}

	// Schedule the workflow.
	if err := r.s.AddWorkflow(ctx, w, jobs, volumes); err != nil {
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

	err = r.p.DeleteVolumesByWorkflowID(ctx, w.ID)
	if err != nil {
		return nil, err
	}

	return &WorkflowResolver{w, r.p}, nil
}

func createVolumes(ctx context.Context, p Persister, w model.Workflow, specs *[]volumeSpec) (map[string]model.Volume, error) {
	volumes := make(map[string]model.Volume)
	if specs != nil {
		for _, vs := range *specs {
			if _, ok := volumes[vs.Name]; ok {
				return nil, fmt.Errorf("duplicate volume declarations")
			}

			v, err := p.CreateVolume(ctx, model.Volume{
				WorkflowID: w.ID,
				Name:       vs.Name,
				Type:       vs.Type,
			})
			if err != nil {
				return nil, err
			}

			volumes[v.Name] = v
		}
	}

	return volumes, nil
}

func createJobs(ctx context.Context, p Persister, w model.Workflow, volumes map[string]model.Volume, specs []jobSpec) ([]model.Job, error) {
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
	var jobs []model.Job
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
		volumeReqs := []model.VolumeRequirement{}
		if js.Volumes != nil {
			for _, v := range *js.Volumes {
				volumeReqs = append(volumeReqs, model.VolumeRequirement{
					Name:     v.Name,
					Location: v.Location,
					VolumeID: volumes[v.Name].ID,
				})
			}
		}

		j, err := p.CreateJob(ctx, model.Job{
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
