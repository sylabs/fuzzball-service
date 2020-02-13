// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package mongodb

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const volumeCollectionName = "volumes"

// CreateVolume creates a new volume. If an ID is provided in v, it is ignored and replaced
// with a unique identifier in the returned volume.
func (c *Connection) CreateVolume(ctx context.Context, v model.Volume) (model.Volume, error) {
	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	v.ID = ""
	ir, err := c.db.Collection(volumeCollectionName).InsertOne(ctx, v)
	if err != nil {
		return model.Volume{}, fmt.Errorf("failed to create volume: %w", err)
	}

	// We want the DB cluster to generate an ID, to ensure it's globally unique.
	v.ID = ir.InsertedID.(primitive.ObjectID).Hex()
	return v, nil
}

// deleteVolume deletes a volume by ID. If the supplied ID is not valid, or there there is not
// a volume with a matching ID in the database, an error is returned.
func (c *Connection) deleteVolume(ctx context.Context, id string) (v model.Volume, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Volume{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(volumeCollectionName).FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&v)
	if err != nil {
		return model.Volume{}, fmt.Errorf("failed to delete volume: %w", err)
	}
	return v, nil
}

// DeleteVolumesByWorkflowID deletes volumes with the given workflow ID.
func (c *Connection) DeleteVolumesByWorkflowID(ctx context.Context, wid string) error {
	_, err := c.db.Collection(volumeCollectionName).DeleteMany(ctx, bson.M{"workflowID": wid})
	if err != nil {
		return fmt.Errorf("failed to delete volumes: %w", err)
	}
	return nil
}

// getVolume retrieves a volume by ID. If the supplied ID is not valid, or there there is not a
// volume with a matching ID in the database, an error is returned.
func (c *Connection) getVolume(ctx context.Context, id string) (v model.Volume, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Volume{}, fmt.Errorf("failed to convert object ID: %w", err)
	}
	err = c.db.Collection(volumeCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&v)
	if err != nil {
		return model.Volume{}, fmt.Errorf("failed to get volume: %w", err)
	}
	return v, nil
}

// GetVolumes returns a list of all volumes.
func (c *Connection) GetVolumes(ctx context.Context, pa model.PageArgs) (p model.VolumesPage, err error) {
	pi, tc, err := findPageEx(ctx, c.db.Collection(volumeCollectionName), maxPageSize, bson.M{}, pa, &p.Volumes)
	if err != nil {
		return p, err
	}
	p.PageInfo = pi
	p.TotalCount = tc
	return p, nil
}

// GetVolumesByWorkflowID returns a list of all volumes required for a given workflow.
func (c *Connection) GetVolumesByWorkflowID(ctx context.Context, pa model.PageArgs, wid string) (p model.VolumesPage, err error) {
	pi, tc, err := findPageEx(ctx, c.db.Collection(volumeCollectionName), maxPageSize, bson.M{"workflowID": wid}, pa, &p.Volumes)
	if err != nil {
		return p, err
	}
	p.PageInfo = pi
	p.TotalCount = tc
	return p, nil
}
