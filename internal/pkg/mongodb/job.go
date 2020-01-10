// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const jobCollectionName = "jobs"

// CreateJob creates a new job.
func (c *Connection) CreateJob(ctx context.Context, j model.Job) (model.Job, error) {
	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	j.ID = ""
	ir, err := c.db.Collection(jobCollectionName).InsertOne(ctx, j)
	if err != nil {
		return model.Job{}, err
	}
	j.ID = ir.InsertedID.(primitive.ObjectID).Hex()
	return j, nil
}

// GetJob retrieves a job by ID.
func (c *Connection) GetJob(ctx context.Context, id string) (j model.Job, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Job{}, err
	}
	err = c.db.Collection(jobCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&j)
	if err != nil {
		return model.Job{}, err
	}
	return j, nil
}
