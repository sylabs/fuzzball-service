// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"
	"fmt"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

// volumeSpec represents a volume specification
type volumeSpec struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

// VolumePersister is the interface by which workflows are persisted.
type VolumePersister interface {
	CreateVolume(context.Context, model.Volume) (model.Volume, error)
	DeleteVolumesByWorkflowID(context.Context, string) error
	GetVolumes(context.Context, model.PageArgs) (model.VolumesPage, error)
	GetVolumesByWorkflowID(context.Context, model.PageArgs, string) (model.VolumesPage, error)
}

// CreateVolume creates a new volume. If an ID is provided in v, it is ignored and replaced
// with a unique identifier in the returned volume.
func (c *Core) CreateVolume(ctx context.Context, v model.Volume) (model.Volume, error) {
	return c.p.CreateVolume(ctx, v)
}

// DeleteVolumesByWorkflowID deletes volumes with the given workflow ID.
func (c *Core) DeleteVolumesByWorkflowID(ctx context.Context, wid string) error {
	return c.p.DeleteVolumesByWorkflowID(ctx, wid)
}

// GetVolumes returns a list of all volumes.
func (c *Core) GetVolumes(ctx context.Context, pa model.PageArgs) (p model.VolumesPage, err error) {
	return c.p.GetVolumes(ctx, pa)
}

// GetVolumesByWorkflowID returns a list of all volumes required for a given workflow.
func (c *Core) GetVolumesByWorkflowID(ctx context.Context, pa model.PageArgs, wid string) (p model.VolumesPage, err error) {
	return c.p.GetVolumesByWorkflowID(ctx, pa, wid)
}

func createVolumes(ctx context.Context, p Persister, w model.Workflow, specs *[]volumeSpec) (map[string]model.Volume, error) {
	volumes := make(map[string]model.Volume)
	if specs != nil {
		for _, vs := range *specs {
			if _, ok := volumes[vs.Name]; ok {
				return nil, fmt.Errorf("duplicate volume declarations")
			}

			v, err := p.CreateVolume(ctx, model.Volume{
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
