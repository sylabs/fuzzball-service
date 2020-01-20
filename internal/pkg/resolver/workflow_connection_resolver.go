// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import "github.com/sylabs/compute-service/internal/pkg/model"

// WorkflowEdgeResolver resolves a workflow edge.
type WorkflowEdgeResolver struct {
	w model.Workflow
}

// Cursor resolves a cursor for use in pagination.
func (r *WorkflowEdgeResolver) Cursor() string {
	return r.w.ID
}

// Node resolves the item at the end of the edge.
func (r *WorkflowEdgeResolver) Node() *WorkflowResolver {
	return &WorkflowResolver{r.w}
}

// WorkflowConnectionResolver resolves a workflow connection.
type WorkflowConnectionResolver struct {
	p model.WorkflowsPage
}

// Edges resolves a list of edges.
func (r *WorkflowConnectionResolver) Edges() *[]*WorkflowEdgeResolver {
	wer := []*WorkflowEdgeResolver{}
	for _, w := range r.p.Workflows {
		wer = append(wer, &WorkflowEdgeResolver{w})
	}
	return &wer
}

// PageInfo resolves information to aid in pagination.
func (r *WorkflowConnectionResolver) PageInfo() *PageInfoResolver {
	return &PageInfoResolver{r.p.PageInfo}
}

// TotalCount resolves the total count of items in the connection.
func (r *WorkflowConnectionResolver) TotalCount() int32 {
	return int32(r.p.TotalCount)
}
