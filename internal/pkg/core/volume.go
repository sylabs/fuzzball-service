// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"fmt"
)

// Volume describes a storage volume.
type Volume struct {
	ID         string `bson:"_id,omitempty"`
	WorkflowID string `bson:"workflowID"`
	Name       string `bson:"name"`
	Type       string `bson:"type"`
}

// VolumesPage represents a page of Volumes resulting from a query, and associated metadata.
type VolumesPage struct {
	Volumes    []Volume // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
}

// volumeSpec represents a volume specification
type volumeSpec struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

// VolumePersister is the interface by which workflows are persisted.
type VolumePersister interface {
	CreateVolume(context.Context, Volume) (Volume, error)
	DeleteVolumesByWorkflowID(context.Context, string) error
	GetVolumes(context.Context, PageArgs) (VolumesPage, error)
	GetVolumesByWorkflowID(context.Context, PageArgs, string) (VolumesPage, error)
}

// CreateVolume creates a new volume. If an ID is provided in v, it is ignored and replaced
// with a unique identifier in the returned volume.
func (c *Core) CreateVolume(ctx context.Context, v Volume) (Volume, error) {
	return c.p.CreateVolume(ctx, v)
}

// DeleteVolumesByWorkflowID deletes volumes with the given workflow ID.
func (c *Core) DeleteVolumesByWorkflowID(ctx context.Context, wid string) error {
	return c.p.DeleteVolumesByWorkflowID(ctx, wid)
}

// GetVolumes returns a list of all volumes.
func (c *Core) GetVolumes(ctx context.Context, pa PageArgs) (p VolumesPage, err error) {
	return c.p.GetVolumes(ctx, pa)
}

// GetVolumesByWorkflowID returns a list of all volumes required for a given workflow.
func (c *Core) GetVolumesByWorkflowID(ctx context.Context, pa PageArgs, wid string) (p VolumesPage, err error) {
	return c.p.GetVolumesByWorkflowID(ctx, pa, wid)
}

func createVolumes(ctx context.Context, p Persister, w Workflow, specs *[]volumeSpec) (map[string]Volume, error) {
	volumes := make(map[string]Volume)
	if specs != nil {
		for _, vs := range *specs {
			if _, ok := volumes[vs.Name]; ok {
				return nil, fmt.Errorf("duplicate volume declarations")
			}

			v, err := p.CreateVolume(ctx, Volume{
				WorkflowID: w.ID,
				Name:       vs.Name,
				Type:       vs.Type,
			})
			if err != nil {
				return nil, err
			}

			volumes[v.Name] = v
		}
	}

	return volumes, nil
}
