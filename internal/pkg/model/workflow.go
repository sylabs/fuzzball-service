// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package model

// Workflow represents a workflow.
type Workflow struct {
	ID     string `bson:"_id,omitempty"`
	Name   string `bson:"name"`
	Status string `bson:"status"`
}

// WorkflowsPage represents a page of workflows resulting from a query, and associated metadata.
type WorkflowsPage struct {
	Workflows  []Workflow // Slice of results.
	PageInfo   PageInfo   // Information to aid in pagination.
	TotalCount int        // Identifies the total count of items in the connection.
}
