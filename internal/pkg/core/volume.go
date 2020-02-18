// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package core

import (
	"context"

	"github.com/sylabs/compute-service/internal/pkg/model"
)

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
