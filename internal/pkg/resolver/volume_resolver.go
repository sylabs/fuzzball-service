package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/sylabs/compute-service/internal/pkg/core"
)

// VolumeServicer is the interface by which volumes are serviced.
type VolumeServicer interface {
	GetVolumes(context.Context, core.PageArgs) (core.VolumesPage, error)
	GetVolumesByWorkflowID(context.Context, core.PageArgs, string) (core.VolumesPage, error)
}

// VolumeResolver resolves a volume.
type VolumeResolver struct {
	v core.Volume
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
