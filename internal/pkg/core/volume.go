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

	c *Core // Used internally for lazy loading.
}

// VolumesPage represents a page of Volumes resulting from a query, and associated metadata.
type VolumesPage struct {
	Volumes    []Volume // Slice of results.
	PageInfo   PageInfo // Information to aid in pagination.
	TotalCount int      // Identifies the total count of items in the connection.
}

// setCore sets the core field of each volume in page p to c.
func (p *VolumesPage) setCore(c *Core) {
	for i := range p.Volumes {
		p.Volumes[i].setCore(c)
	}
}

// setCore sets the core of v to c.
func (v *Volume) setCore(c *Core) {
	v.c = c
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
