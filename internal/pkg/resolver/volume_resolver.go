package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// VolumePersister is the interface by which workflows are persisted.
type VolumePersister interface {
	CreateVolume(context.Context, model.Volume) (model.Volume, error)
	DeleteVolumesByWorkflowID(context.Context, string) error
	GetVolumes(context.Context, model.PageArgs) (model.VolumesPage, error)
	GetVolumesByWorkflowID(context.Context, model.PageArgs, string) (model.VolumesPage, error)
}

// volumeSpec represents a volume specification
type volumeSpec struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

// VolumeResolver resolves a volume.
type VolumeResolver struct {
	v model.Volume
	p Persister
}

// ID resolves the volume ID.
func (r *VolumeResolver) ID() graphql.ID {
	return graphql.ID(r.v.ID)
}

// Name resolves the volume name.
func (r *VolumeResolver) Name() string {
	return r.v.Name
}

// Type resolves the volume type.
func (r *VolumeResolver) Type() string {
	return r.v.Type
}
