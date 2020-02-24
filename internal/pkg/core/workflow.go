// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"fmt"
	"time"

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
	ID         string     `bson:"_id,omitempty"`
	CreatedAt  time.Time  `bson:"createdAt"`
	StartedAt  *time.Time `bson:"startedAt,omitempty"`
	FinishedAt *time.Time `bson:"finishedAt,omitempty"`
	Name       string     `bson:"name"`
	Status     string     `bson:"status"`

	c *Core // Used internally for lazy loading.
}

// CreatedBy retrieves the user that created workflow w.
func (w Workflow) CreatedBy(ctx context.Context) (User, error) {
	u := User{
		ID:    "507f1f77bcf86cd799439011",
		Login: "jimbob",
	}
	u.setCore(w.c)
	return u, nil
}

// JobsPage retrieves a page of jobs related to workflow w.
func (w Workflow) JobsPage(ctx context.Context, pa PageArgs) (JobsPage, error) {
	p, err := w.c.p.GetJobsByWorkflowID(ctx, pa, w.ID)
	if err != nil {
		return JobsPage{}, err
	}
	p.setCore(w.c)
	return p, nil
}

// VolumesPage retrieves a page of volumes related to workflow w.
func (w Workflow) VolumesPage(ctx context.Context, pa PageArgs) (VolumesPage, error) {
	p, err := w.c.p.GetVolumesByWorkflowID(ctx, pa, w.ID)
	if err != nil {
		return VolumesPage{}, err
	}
	p.setCore(w.c)
	return p, nil
}

// setCore sets the core of w to c.
func (w *Workflow) setCore(c *Core) {
	w.c = c
}

// WorkflowsPage represents a page of workflows resulting from a query, and associated metadata.
type WorkflowsPage struct {
	Workflows  []Workflow // Slice of results.
	PageInfo   PageInfo   // Information to aid in pagination.
	TotalCount int        // Identifies the total count of items in the connection.
}

// setCore sets the core field of each workflow in page p to c.
func (p *WorkflowsPage) setCore(c *Core) {
	for i := range p.Workflows {
		p.Workflows[i].setCore(c)
	}
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

		j, err := c.p.CreateJob(ctx, Job{
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
