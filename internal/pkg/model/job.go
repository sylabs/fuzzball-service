// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package model

type Job struct {
	ID         string              `bson:"_id,omitempty"`
	WorkflowID string              `bson:"workflowID"`
	Name       string              `bson:"name"`
	Image      string              `bson:"image"`
	Command    []string            `bson:"command"`
	Status     string              `bson:"status"`
	ExitCode   int                 `bson:"exitCode"`
	Requires   []string            `bson:"requires"`
	Volumes    []VolumeRequirement `bson:"volumes"`
}

type VolumeRequirement struct {
	VolumeID string `bson:"volumeID"`
	Name     string `bson:"name"`
	Location string `bson:"location"`
}

// JobsPage represents a page of jobs resulting from a query, and associated metadata.
type JobsPage struct {
	Jobs       []Job    // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
}
