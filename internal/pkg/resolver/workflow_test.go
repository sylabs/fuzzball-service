// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/sylabs/compute-service/internal/pkg/core"
	"github.com/sylabs/compute-service/internal/pkg/schema"
)

func TestWorkflow(t *testing.T) {
	mc, err := getMockCore(mockCore{
		p: mockPersister{
			w: core.Workflow{
				ID:        "workflowID",
				Name:      "workflowName",
				CreatedAt: time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := schema.Get(&Resolver{s: mc})
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
	mc, err := getMockCore(mockCore{
		p: mockPersister{
			w: core.Workflow{
				ID:        "workflowID",
				Name:      "workflowName",
				CreatedAt: time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC),
			},
			j: core.Job{
				ID:         "jobID",
				WorkflowID: "workflowID",
				Name:       "jobName",
				Image:      "jobImage",
				Command:    []string{"jobCommand"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := schema.Get(&Resolver{s: mc})
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
	mc, err := getMockCore(mockCore{
		p: mockPersister{
			w: core.Workflow{
				ID:        "workflowID",
				Name:      "workflowName",
				CreatedAt: time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC),
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
	})
	if err != nil {
		t.Fatal(err)
	}

	s, err := schema.Get(&Resolver{s: mc})
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
