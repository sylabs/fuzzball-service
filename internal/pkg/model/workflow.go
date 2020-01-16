// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package model

// Workflow represents a workflow.
type Workflow struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
}
