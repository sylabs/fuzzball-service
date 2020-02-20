// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/graph"
)

// Job contains information about an indivisual job.
type Job struct {
	ID         string              `bson:"_id,omitempty"`
	WorkflowID string              `bson:"workflowID"`
	Name       string              `bson:"name"`
	Image      string              `bson:"image"`
	Command    []string            `bson:"command"`
	Status     string              `bson:"status"`
	ExitCode   int                 `bson:"exitCode"`
	Requires   []string            `bson:"requires"`
	Volumes    []VolumeRequirement `bson:"volumes"`
}

// VolumeRequirement describes a required volume.
type VolumeRequirement struct {
	VolumeID string `bson:"volumeID"`
	Name     string `bson:"name"`
	Location string `bson:"location"`
}

// JobsPage represents a page of jobs resulting from a query, and associated metadata.
type JobsPage struct {
	Jobs       []Job    // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
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

// JobPersister is the interface by which jobs are persisted.
type JobPersister interface {
	CreateJob(context.Context, Job) (Job, error)
	DeleteJobsByWorkflowID(context.Context, string) error
	GetJob(context.Context, string) (Job, error)
	GetJobs(context.Context, PageArgs) (JobsPage, error)
	GetJobsByWorkflowID(context.Context, PageArgs, string) (JobsPage, error)
	GetJobsByID(context.Context, PageArgs, string, []string) (JobsPage, error)
}

// CreateJob creates a new job. If an ID is provided in j, it is ignored and replaced with a unique
// identifier in the returned job.
func (c *Core) CreateJob(ctx context.Context, j Job) (Job, error) {
	return c.p.CreateJob(ctx, j)
}

// DeleteJobsByWorkflowID deletes jobs with the given workflow ID.
func (c *Core) DeleteJobsByWorkflowID(ctx context.Context, wid string) error {
	return c.p.DeleteJobsByWorkflowID(ctx, wid)
}

// GetJob retrieves a job by ID. If the supplied ID is not valid, or there there is not a
// job with a matching ID in the database, an error is returned.
func (c *Core) GetJob(ctx context.Context, id string) (j Job, err error) {
	return c.p.GetJob(ctx, id)
}

// GetJobs returns a list of all jobs.
func (c *Core) GetJobs(ctx context.Context, pa PageArgs) (p JobsPage, err error) {
	return c.p.GetJobs(ctx, pa)
}

// GetJobsByWorkflowID returns a list of all jobs for a given workflow.
func (c *Core) GetJobsByWorkflowID(ctx context.Context, pa PageArgs, wid string) (p JobsPage, err error) {
	return c.p.GetJobsByWorkflowID(ctx, pa, wid)
}

// GetJobsByID returns a list of jobs by name within a given workflow.
func (c *Core) GetJobsByID(ctx context.Context, pa PageArgs, wid string, ids []string) (p JobsPage, err error) {
	return c.p.GetJobsByID(ctx, pa, wid, ids)
}

func createJobs(ctx context.Context, p Persister, w Workflow, volumes map[string]Volume, specs []jobSpec) ([]Job, error) {
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

		j, err := p.CreateJob(ctx, Job{
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
