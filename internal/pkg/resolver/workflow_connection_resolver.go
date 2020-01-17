// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import "github.com/sylabs/compute-service/internal/pkg/model"

// WorkflowEdgeResolver resolves a workflow edge.
type WorkflowEdgeResolver struct {
}

// Cursor resolves a cursor for use in pagination.
func (r *WorkflowEdgeResolver) Cursor() string {
	return "" // TODO
}

// Node resolves the item at the end of the edge.
func (r *WorkflowEdgeResolver) Node() *WorkflowResolver {
	return nil // TODO
}

// WorkflowConnectionResolver resolves a workflow connection.
type WorkflowConnectionResolver struct {
	p model.WorkflowsPage
}

// Edges resolves a list of edges.
func (r *WorkflowConnectionResolver) Edges() *[]*WorkflowEdgeResolver {
	return nil // TODO
}

// Nodes resolves a list of nodes.
func (r *WorkflowConnectionResolver) Nodes() *[]*WorkflowResolver {
	wr := []*WorkflowResolver{}
	for _, w := range r.p.Workflows {
		wr = append(wr, &WorkflowResolver{&w})
	}
	return &wr
}

// PageInfo resolves information to aid in pagination.
func (r *WorkflowConnectionResolver) PageInfo() *PageInfoResolver {
	return &PageInfoResolver{r.p.PageInfo}
}

// TotalCount resolves the total count of items in the connection.
func (r *WorkflowConnectionResolver) TotalCount() int32 {
	return int32(r.p.TotalCount)
}
