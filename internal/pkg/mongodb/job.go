// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const jobCollectionName = "jobs"

// CreateJob creates a new job. If an ID is provided in j, it is ignored and replaced with a unique
// identifier in the returned job.
func (c *Connection) CreateJob(ctx context.Context, j model.Job) (model.Job, error) {
	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	j.ID = ""
	ir, err := c.db.Collection(jobCollectionName).InsertOne(ctx, j)
	if err != nil {
		return model.Job{}, fmt.Errorf("failed to create job: %w", err)
	}
	j.ID = ir.InsertedID.(primitive.ObjectID).Hex()
	return j, nil
}

// DeleteJob deletes a job by ID. If the supplied ID is not valid, or there there is not a job with
// a matching ID in the database, an error is returned.
func (c *Connection) DeleteJob(ctx context.Context, id string) (j model.Job, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Job{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(jobCollectionName).FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&j)
	if err != nil {
		return model.Job{}, fmt.Errorf("failed to delete job: %w", err)
	}
	return j, nil
}

// GetJobs returns all jobs if no search arguments specified, otherwise
// returns matching job
func (c *Connection) GetJobs(ctx context.Context, filterSpec map[string]string) (j []model.Job, err error) {
	var oid primitive.ObjectID

	filterID, ok := filterSpec["ID"]
	if ok && filterID != "" {
		// Handle 'ID'
		var err error
		oid, err = primitive.ObjectIDFromHex(filterID)
		if err != nil {
			return nil, fmt.Errorf("failed to convert object ID: %w", err)
		}
	}

	var s bson.M

	if oid != primitive.NilObjectID {
		s = bson.M{"_id": oid}
	} else {
		s = bson.M{}
		for key, value := range filterSpec {
			s[key] = value
		}
	}

	cur, err := c.db.Collection(jobCollectionName).Find(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("failed to get job(s): %w", err)
	}

	var results []model.Job
	for cur.Next(ctx) {
		var j model.Job
		if err := cur.Decode(&j); err != nil {
			continue
		}
		results = append(results, j)
	}
	return results, nil
}
