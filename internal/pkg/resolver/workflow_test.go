// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"testing"
	"time"

	"github.com/sylabs/compute-service/internal/pkg/core"
	"github.com/sylabs/compute-service/internal/pkg/schema"
)

func TestWorkflowJobs(t *testing.T) {
	startedAt := time.Date(2020, 01, 20, 19, 21, 31, 0, time.UTC)
	finishedAt := time.Date(2020, 01, 20, 19, 21, 32, 0, time.UTC)
	w := core.Workflow{
		ID:         "workflowID",
		CreatedAt:  time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC),
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
		Name:       "workflowName",
		Status:     "COMPLETED",
	}

	sc := "startCursor"
	ec := "endCursor"
	jp := core.JobsPage{
		Jobs: []core.Job{
			{
				ID:   "id1",
				Name: "name1",
			},
			{
				ID:   "id2",
				Name: "name2",
			},
		},
		PageInfo: core.PageInfo{
			StartCursor:     &sc,
			EndCursor:       &ec,
			HasNextPage:     true,
			HasPreviousPage: false,
		},
		TotalCount: 2,
	}

	cursor := "cursorValue"
	count := 2

	tests := []struct {
		name   string
		args   map[string]interface{}
		wantPA core.PageArgs
	}{
		{"OK", map[string]interface{}{"id": "workflowID"}, core.PageArgs{}},
		{"After", map[string]interface{}{"id": "workflowID", "after": cursor}, core.PageArgs{After: &cursor}},
		{"Before", map[string]interface{}{"id": "workflowID", "before": cursor}, core.PageArgs{Before: &cursor}},
		// The first and last params enter as float64s via the HTTP handler, so test that here.
		{"First", map[string]interface{}{"id": "workflowID", "first": float64(count)}, core.PageArgs{First: &count}},
		{"Last", map[string]interface{}{"id": "workflowID", "last": float64(count)}, core.PageArgs{Last: &count}},
		{"BadID", map[string]interface{}{"id": "bad"}, core.PageArgs{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, err := getMockCore(mockCore{
				p: mockPersister{
					wantPA: tt.wantPA,
					w:      w,
					jp:     jp,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			s, err := schema.Get(&Resolver{s: mc})
			if err != nil {
				t.Fatal(err)
			}

			q := `
			query OpName($id: ID!, $after: String, $before: String, $first: Int, $last: Int)
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
			    status
			    jobs(after: $after, before: $before, first: $first, last: $last) {
			      edges {
			        cursor
			        node {
			          id
			          name
			        }
			      }
			      pageInfo {
			        startCursor
			        endCursor
			        hasNextPage
			        hasPreviousPage
			      }
			      totalCount
			    }
			  }
			}`

			res := s.Exec(context.Background(), q, "", tt.args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWorkflowVolumes(t *testing.T) {
	startedAt := time.Date(2020, 01, 20, 19, 21, 31, 0, time.UTC)
	finishedAt := time.Date(2020, 01, 20, 19, 21, 32, 0, time.UTC)
	w := core.Workflow{
		ID:         "workflowID",
		CreatedAt:  time.Date(2020, 01, 20, 19, 21, 30, 0, time.UTC),
		StartedAt:  &startedAt,
		FinishedAt: &finishedAt,
		Name:       "workflowName",
		Status:     "COMPLETED",
	}

	sc := "startCursor"
	ec := "endCursor"
	vp := core.VolumesPage{
		Volumes: []core.Volume{
			{
				ID:   "id1",
				Name: "name1",
			},
			{
				ID:   "id2",
				Name: "name2",
			},
		},
		PageInfo: core.PageInfo{
			StartCursor:     &sc,
			EndCursor:       &ec,
			HasNextPage:     true,
			HasPreviousPage: false,
		},
		TotalCount: 2,
	}

	cursor := "cursorValue"
	count := 2

	tests := []struct {
		name   string
		args   map[string]interface{}
		wantPA core.PageArgs
	}{
		{"OK", map[string]interface{}{"id": "workflowID"}, core.PageArgs{}},
		{"After", map[string]interface{}{"id": "workflowID", "after": cursor}, core.PageArgs{After: &cursor}},
		{"Before", map[string]interface{}{"id": "workflowID", "before": cursor}, core.PageArgs{Before: &cursor}},
		// The first and last params enter as float64s via the HTTP handler, so test that here.
		{"First", map[string]interface{}{"id": "workflowID", "first": float64(count)}, core.PageArgs{First: &count}},
		{"Last", map[string]interface{}{"id": "workflowID", "last": float64(count)}, core.PageArgs{Last: &count}},
		{"BadID", map[string]interface{}{"id": "bad"}, core.PageArgs{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, err := getMockCore(mockCore{
				p: mockPersister{
					wantPA: tt.wantPA,
					w:      w,
					vp:     vp,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			s, err := schema.Get(&Resolver{s: mc})
			if err != nil {
				t.Fatal(err)
			}

			q := `
			query OpName($id: ID!, $after: String, $before: String, $first: Int, $last: Int)
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
			    status
			    volumes(after: $after, before: $before, first: $first, last: $last) {
			      edges {
			        cursor
			        node {
			          id
			          name
			        }
			      }
			      pageInfo {
			        startCursor
			        endCursor
			        hasNextPage
			        hasPreviousPage
			      }
			      totalCount
			    }
			  }
			}`

			res := s.Exec(context.Background(), q, "", tt.args)

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
