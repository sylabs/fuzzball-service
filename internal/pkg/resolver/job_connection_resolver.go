// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import "github.com/sylabs/fuzzball-service/internal/pkg/core"

// JobEdgeResolver resolves a job edge.
type JobEdgeResolver struct {
	j core.Job
}

// Cursor resolves a cursor for use in pagination.
func (r *JobEdgeResolver) Cursor() string {
	return r.j.ID
}

// Node resolves the item at the end of the edge.
func (r *JobEdgeResolver) Node() *JobResolver {
	return &JobResolver{r.j}
}

// JobConnectionResolver resolves a job connection.
type JobConnectionResolver struct {
	jp core.JobsPage
}

// Edges resolves a list of edges.
func (r *JobConnectionResolver) Edges() *[]*JobEdgeResolver {
	wer := []*JobEdgeResolver{}
	for _, w := range r.jp.Jobs {
		wer = append(wer, &JobEdgeResolver{w})
	}
	return &wer
}

// PageInfo resolves information to aid in pagination.
func (r *JobConnectionResolver) PageInfo() *PageInfoResolver {
	return &PageInfoResolver{r.jp.PageInfo}
}

// TotalCount resolves the total count of items in the connection.
func (r *JobConnectionResolver) TotalCount() int32 {
	return int32(r.jp.TotalCount)
}
