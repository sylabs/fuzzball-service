// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package model

type Volume struct {
	ID         string `bson:"_id,omitempty"`
	WorkflowID string `bson:"workflowID"`
	Name       string `bson:"name"`
	Type       string `bson:"type"`
}

// VolumesPage represents a page of Volumes resulting from a query, and associated metadata.
type VolumesPage struct {
	Volumes    []Volume // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
}
