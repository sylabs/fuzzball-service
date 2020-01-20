// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
)

type SystemInfo struct {
	HostName        string         `json:"hostname"`
	CPUArchitecture string         `json:"cpuarchitecture"`
	OSPlatform      string         `json:"osplatform"`
	Memory          uint64         `json:"memory"`
	Capabilities    *[]*Capability `json:"capabilities"`
}

type SystemInfoResolver struct {
	s SystemInfo
}

func (s *SystemInfoResolver) HostName() string {
	return s.s.HostName
}

func (s *SystemInfoResolver) CPUArchitecture() string {
	return s.s.CPUArchitecture
}

func (s *SystemInfoResolver) OSPlatform() string {
	return s.s.OSPlatform
}

func (s *SystemInfoResolver) Memory() int32 {
	return int32(s.s.Memory)
}

func (s *SystemInfoResolver) Capabilities(ctx context.Context) (*[]*CapabilityResolver, error) {
	return NewCapabilities(ctx)
}
