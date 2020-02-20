// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

// PageArgs contains criteria to select a page of results.
type PageArgs struct {
	After  *string // Select elements in the list that come after the specified cursor.
	Before *string // Select elements in the list that come before the specified cursor.
	First  *int    // Select the first n elements from the list.
	Last   *int    // Select the last n elements from the list.
}

// PageInfo contains information to aid in pagination.
type PageInfo struct {
	StartCursor     *string // When paginating backwards, the cursor to continue.
	EndCursor       *string // When paginating forwards, the cursor to continue.
	HasNextPage     bool    // When paginating forwards, are there more items?
	HasPreviousPage bool    // When paginating backwards, are there more items?
}
