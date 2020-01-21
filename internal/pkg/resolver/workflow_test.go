// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

type mockPersister struct {
	wantPA model.PageArgs
	w      model.Workflow
	wp     model.WorkflowsPage
}

func (p *mockPersister) CreateWorkflow(ctx context.Context, w model.Workflow) (model.Workflow, error) {
	if got, want := w.Name, p.w.Name; got != want {
		return model.Workflow{}, fmt.Errorf("got name %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) DeleteWorkflow(ctx context.Context, id string) (model.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return model.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) GetWorkflow(ctx context.Context, id string) (model.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return model.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, nil
}

func (p *mockPersister) GetWorkflows(ctx context.Context, pa model.PageArgs) (model.WorkflowsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return model.WorkflowsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.wp, nil
}

func TestWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: model.Workflow{
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
			w: model.Workflow{
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
		name   string
		wfName string
	}{
		{"OK", "workflowName"},
		{"BadName", "bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := `
			mutation OpName($name: String!) {
			  createWorkflow(name: $name) {
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
				"name": tt.wfName,
			}

			res := s.Exec(context.Background(), q, "", args)

			if err := verifyGoldenJSON(t.Name(), res); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDeleteWorkflow(t *testing.T) {
	r := Resolver{
		p: &mockPersister{
			w: model.Workflow{
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
