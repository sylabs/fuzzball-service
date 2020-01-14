// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

type SystemInfo struct {
	HostName     string         `json:"hostname"`
	Capabilities *[]*Capability `json:"capabilities"`
}

type SystemInfoResolver struct {
	s *SystemInfo
}

func (s *SystemInfoResolver) HostName() string {
	return s.s.HostName
}

func (s *SystemInfoResolver) Capabilities(ctx context.Context) (*[]*CapabilityResolver, error) {
	return NewCapabilities(ctx)
}