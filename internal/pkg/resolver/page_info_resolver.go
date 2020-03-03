// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"github.com/sylabs/fuzzball-service/internal/pkg/core"
)

// PageInfoResolver resolves information about pagination in a connection.
type PageInfoResolver struct {
	pi core.PageInfo
}

// StartCursor resolves the cursor to continue when paginating backwards.
func (r *PageInfoResolver) StartCursor() *string {
	return r.pi.StartCursor
}

// EndCursor resolves the cursor to continue when paginating forwards.
func (r *PageInfoResolver) EndCursor() *string {
	return r.pi.EndCursor
}

// HasNextPage resolves whether there are more items when paginating forwards.
func (r *PageInfoResolver) HasNextPage() bool {
	return r.pi.HasNextPage
}

// HasPreviousPage resolves whether there are more items when paginating backwards.
func (r *PageInfoResolver) HasPreviousPage() bool {
	return r.pi.HasPreviousPage
}
