// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

type Capability struct {
	key string
	value string
}

type CapabilityResolver struct {
	capability *Capability
}

func (r *CapabilityResolver) Key() string {
	return r.capability.key
}

func (r *CapabilityResolver) Value() string {
	return r.capability.value
}
