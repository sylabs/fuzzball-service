// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/core"
)

type mockPersister struct {
	wantPA core.PageArgs
	j      core.Job
	v      core.Volume
	w      core.Workflow
	jp     core.JobsPage
	vp     core.VolumesPage
	wp     core.WorkflowsPage
}

func (p *mockPersister) CreateWorkflow(ctx context.Context, s core.WorkflowSpec) (core.Workflow, error) {
	if got, want := s.Name, p.w.Name; got != want {
		return core.Workflow{}, fmt.Errorf("got name %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) DeleteWorkflow(ctx context.Context, id string) (core.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return core.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) GetWorkflow(ctx context.Context, id string) (core.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return core.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) GetWorkflows(ctx context.Context, pa core.PageArgs) (core.WorkflowsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.WorkflowsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.wp, nil
}

func (p *mockPersister) GetJob(ctx context.Context, id string) (core.Job, error) {
	if got, want := id, p.j.ID; got != want {
		return core.Job{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.j, nil
}

func (p *mockPersister) GetJobs(ctx context.Context, pa core.PageArgs) (core.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, nil
}

func (p *mockPersister) GetJobsByID(ctx context.Context, pa core.PageArgs, wid string, names []string) (core.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, nil
}

func (p *mockPersister) GetJobsByWorkflowID(ctx context.Context, pa core.PageArgs, id string) (core.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, nil
}

func (p *mockPersister) GetVolumes(ctx context.Context, pa core.PageArgs) (core.VolumesPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.VolumesPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.vp, nil
}

func (p *mockPersister) GetVolumesByWorkflowID(ctx context.Context, pa core.PageArgs, id string) (core.VolumesPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.VolumesPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.vp, nil
}

func TestWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: core.Workflow{
				ID:   "workflowID",
				Name: "workflowName",
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		id   string
	}{
		{"OK", "workflowID"},
		{"BadID", "bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			query OpName($id: ID!)
			{
			  workflow(id: $id) {
			    id
			    name
			    createdBy {
			      id
			      login
			    }
			    createdAt
			    startedAt
			    finishedAt
			  }
			}`
			args := map[string]interface{}{
				"id": tt.id,
			}

			res := s.Exec(context.Background(), q, "", args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: core.Workflow{
				ID:   "workflowID",
				Name: "workflowName",
			},
			j: core.Job{
				ID:         "jobID",
				WorkflowID: "workflowID",
				Name:       "jobName",
				Image:      "jobImage",
				Command:    []string{"jobCommand"},
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	okMap := map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "workflowName",
			"jobs": map[string]interface{}{
				"name":    "jobName",
				"image":   "jobImage",
				"command": "jobCommand",
			},
		},
	}

	badMap := map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "bad",
			"jobs": map[string]interface{}{
				"name":    "jobName",
				"image":   "jobImage",
				"command": "jobCommand",
			},
		},
	}

	tests := []struct {
		name string
		vars map[string]interface{}
	}{
		{"OK", okMap},
		{"BadName", badMap},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			mutation OpName($spec: WorkflowSpec!) {
			  createWorkflow(spec: $spec) {
			    id
			    name
			    createdBy {
			      id
			      login
			    }
			    createdAt
			    startedAt
			    finishedAt
			  }
			}`

			res := s.Exec(context.Background(), q, "", tt.vars)
			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: core.Workflow{
				ID:   "workflowID",
				Name: "workflowName",
			},
			j: core.Job{
				ID:         "jobID",
				WorkflowID: "workflowID",
				Name:       "jobName",
				Image:      "jobImage",
				Command:    []string{"jobCommand"},
			},
			v: core.Volume{
				ID:         "volumeID",
				WorkflowID: "workflowID",
				Name:       "volumeName",
				Type:       "volumeType",
			},
		},
	}
	s, err := getSchema(&r)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		id   string
	}{
		{"OK", "workflowID"},
		{"BadID", "bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			mutation OpName($id: ID!) {
			  deleteWorkflow(id: $id) {
			    id
			    name
			    createdBy {
			      id
			      login
			    }
			    createdAt
			    startedAt
			    finishedAt
			  }
			}`

			args := map[string]interface{}{
				"id": tt.id,
			}

			res := s.Exec(context.Background(), q, "", args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}
