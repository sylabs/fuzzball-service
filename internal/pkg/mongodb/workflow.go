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
func (c *Connection) CreateWorkflow(ctx context.Context, w model.Workflow) (model.Workflow, error) {
	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	w.ID = ""
	ir, err := c.db.Collection(workflowCollectionName).InsertOne(ctx, w)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to create workflow: %w", err)
	}
	w.ID = ir.InsertedID.(primitive.ObjectID).Hex()
	return w, nil
}

// DeleteWorkflow deletes a workflow by ID. If the supplied ID is not valid, or there there is not
// a workflow with a matching ID in the database, an error is returned.
func (c *Connection) DeleteWorkflow(ctx context.Context, id string) (w model.Workflow, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(workflowCollectionName).FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&w)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to delete workflow: %w", err)
	}
	return w, nil
}

// GetWorkflow retrieves a workflow by ID. If the supplied ID is not valid, or there there is not a
// workflow with a matching ID in the database, an error is returned.
func (c *Connection) GetWorkflow(ctx context.Context, id string) (w model.Workflow, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(workflowCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&w)
	if err != nil {
		return model.Workflow{}, fmt.Errorf("failed to get workflow: %w", err)
	}
	return w, nil
}

// GetWorkflows returns a list of all workflows.
//
// TODO: pagination
func (c *Connection) GetWorkflows(ctx context.Context, pa model.PageArgs) (p model.WorkflowsPage, err error) {
	cur, err := c.db.Collection(workflowCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return p, fmt.Errorf("failed to get workflow: %w", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var w model.Workflow
		if err := cur.Decode(&w); err != nil {
			return p, fmt.Errorf("failed to decode workflow: %w", err)
		}
		p.Workflows = append(p.Workflows, w)
		p.TotalCount++
	}
	if err := cur.Err(); err != nil {
		return p, fmt.Errorf("failed to get workflows: %w", err)
	}
	return p, nil
}
