// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const jobCollectionName = "jobs"

// CreateJob creates a new job. If an ID is provided in j, it is ignored and replaced
// with a unique identifier in the returned job.
func (c *Connection) CreateJob(ctx context.Context, j core.Job) (core.Job, error) {
	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	j.ID = ""
	ir, err := c.db.Collection(jobCollectionName).InsertOne(ctx, j)
	if err != nil {
		return core.Job{}, fmt.Errorf("failed to create job: %w", err)
	}

	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	j.ID = ir.InsertedID.(primitive.ObjectID).Hex()
	return j, nil
}

// DeleteJobsByWorkflowID deletes jobs with the given workflow ID.
func (c *Connection) DeleteJobsByWorkflowID(ctx context.Context, wid string) error {
	_, err := c.db.Collection(jobCollectionName).DeleteMany(ctx, bson.M{"workflowID": wid})
	if err != nil {
		return fmt.Errorf("failed to delete jobs: %w", err)
	}
	return nil
}

// deleteJob deletes a job by ID. If the supplied ID is not valid, or there there is not
// a job with a matching ID in the database, an error is returned.
func (c *Connection) deleteJob(ctx context.Context, id string) (j core.Job, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return core.Job{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(jobCollectionName).FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&j)
	if err != nil {
		return core.Job{}, fmt.Errorf("failed to delete job: %w", err)
	}
	return j, nil
}

// GetJob retrieves a job by ID. If the supplied ID is not valid, or there there is not a
// job with a matching ID in the database, an error is returned.
func (c *Connection) GetJob(ctx context.Context, id string) (j core.Job, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return core.Job{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(jobCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&j)
	if err != nil {
		return core.Job{}, fmt.Errorf("failed to get job: %w", err)
	}
	return j, nil
}

// GetJobs returns a list of all jobs.
func (c *Connection) GetJobs(ctx context.Context, pa core.PageArgs) (p core.JobsPage, err error) {
	pi, tc, err := findPageEx(ctx, c.db.Collection(jobCollectionName), maxPageSize, bson.M{}, pa, &p.Jobs)
	if err != nil {
		return p, err
	}
	p.PageInfo = pi
	p.TotalCount = tc
	return p, nil
}

// GetJobsByID returns a list of jobs by name within a given workflow.
func (c *Connection) GetJobsByID(ctx context.Context, pa core.PageArgs, wid string, ids []string) (p core.JobsPage, err error) {
	// short circuit if we have no ids to look up
	// mongo does not like an empty array passed
	// with the $in parameter
	if len(ids) == 0 {
		return p, nil
	}

	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return p, fmt.Errorf("failed to convert object ID: %w", err)
		}
		oids = append(oids, oid)
	}

	filter := bson.M{"workflowID": wid, "_id": bson.M{"$in": oids}}
	pi, tc, err := findPageEx(ctx, c.db.Collection(jobCollectionName), maxPageSize, filter, pa, &p.Jobs)
	if err != nil {
		return p, err
	}
	p.PageInfo = pi
	p.TotalCount = tc
	return p, nil
}

// GetJobsByWorkflowID returns a list of all jobs for a given workflow.
func (c *Connection) GetJobsByWorkflowID(ctx context.Context, pa core.PageArgs, wid string) (p core.JobsPage, err error) {
	pi, tc, err := findPageEx(ctx, c.db.Collection(jobCollectionName), maxPageSize, bson.M{"workflowID": wid}, pa, &p.Jobs)
	if err != nil {
		return p, err
	}
	p.PageInfo = pi
	p.TotalCount = tc
	return p, nil
}

// SetJobStatus updates a job's status and records the supplied
// exit status. If the supplied ID is not valid, or there there is
// not a job with a matching ID in the database, an error is returned.
// TODO(ian): should probably have the exit status broken out into another method.
func (c *Connection) SetJobStatus(ctx context.Context, id, status string, exitCode int) (err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert object ID: %w", err)
	}

	update := bson.M{"$set": bson.M{"status": status, "exitCode": exitCode}}
	err = c.db.Collection(jobCollectionName).FindOneAndUpdate(ctx, bson.M{"_id": oid}, update).Err()
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}
	return nil
}
