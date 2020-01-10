// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package model

// Job represents a job.
type Job struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
}
