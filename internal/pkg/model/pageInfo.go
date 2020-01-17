// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package model

// PageInfo contains information to aid in pagination.
type PageInfo struct {
	StartCursor     *string // When paginating backwards, the cursor to continue.
	EndCursor       *string // When paginating forwards, the cursor to continue.
	HasNextPage     bool    // When paginating forwards, are there more items?
	HasPreviousPage bool    // When paginating backwards, are there more items?
}
