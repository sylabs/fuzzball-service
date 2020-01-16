// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const workflowCollectionName = "workflows"

// CreateWorkflow creates a new workflow. If an ID is provided in j, it is ignored and replaced
// with a unique identifier in the returned workflow.
func (c *Connection) CreateWorkflow(ctx context.Context, j model.Workflow) (model.Workflow, error) {
	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	j.ID = ""
	ir, err := c.db.Collection(workflowCollectionName).InsertOne(ctx, j)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to create workflow: %w", err)
	}
	j.ID = ir.InsertedID.(primitive.ObjectID).Hex()
	return j, nil
}

// DeleteWorkflow deletes a workflow by ID. If the supplied ID is not valid, or there there is not
// a workflow with a matching ID in the database, an error is returned.
func (c *Connection) DeleteWorkflow(ctx context.Context, id string) (j model.Workflow, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(workflowCollectionName).FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&j)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to delete workflow: %w", err)
	}
	return j, nil
}

// GetWorkflow retrieves a workflow by ID. If the supplied ID is not valid, or there there is not a
// workflow with a matching ID in the database, an error is returned.
func (c *Connection) GetWorkflow(ctx context.Context, id string) (j model.Workflow, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(workflowCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&j)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to get workflow: %w", err)
	}
	return j, nil
}
