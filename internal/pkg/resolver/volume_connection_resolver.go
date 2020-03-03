// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import "github.com/sylabs/fuzzball-service/internal/pkg/core"

// VolumeEdgeResolver resolves a volume edge.
type VolumeEdgeResolver struct {
	v core.Volume
}

// Cursor resolves a cursor for use in pagination.
func (r *VolumeEdgeResolver) Cursor() string {
	return r.v.ID
}

// Node resolves the item at the end of the edge.
func (r *VolumeEdgeResolver) Node() *VolumeResolver {
	return &VolumeResolver{r.v}
}

// VolumeConnectionResolver resolves a volume connection.
type VolumeConnectionResolver struct {
	vp core.VolumesPage
}

// Edges resolves a list of edges.
func (r *VolumeConnectionResolver) Edges() *[]*VolumeEdgeResolver {
	wer := []*VolumeEdgeResolver{}
	for _, w := range r.vp.Volumes {
		wer = append(wer, &VolumeEdgeResolver{w})
	}
	return &wer
}

// PageInfo resolves information to aid in pagination.
func (r *VolumeConnectionResolver) PageInfo() *PageInfoResolver {
	return &PageInfoResolver{r.vp.PageInfo}
}

// TotalCount resolves the total count of items in the connection.
func (r *VolumeConnectionResolver) TotalCount() int32 {
	return int32(r.vp.TotalCount)
}
