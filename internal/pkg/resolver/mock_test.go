// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sylabs/compute-service/internal/pkg/core"
)

type mockPersister struct {
	wantPA core.PageArgs
	j      core.Job
	v      core.Volume
	w      core.Workflow
	jp     core.JobsPage
	vp     core.VolumesPage
	wp     core.WorkflowsPage
	err    error
}

func (p mockPersister) CreateWorkflow(ctx context.Context, w core.Workflow) (core.Workflow, error) {
	if got, want := w.Name, p.w.Name; got != want {
		return core.Workflow{}, fmt.Errorf("got name %v, want %v", got, want)
	}
	return p.w, p.err
}

func (p mockPersister) DeleteWorkflow(ctx context.Context, id string) (core.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return core.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, p.err
}

func (p mockPersister) GetWorkflow(ctx context.Context, id string) (core.Workflow, error) {
	if got, want := id, p.w.ID; got != want {
		return core.Workflow{}, fmt.Errorf("got ID %v, want %v", got, want)
	}
	return p.w, p.err
}

func (p mockPersister) GetWorkflows(ctx context.Context, pa core.PageArgs) (core.WorkflowsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.WorkflowsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.wp, p.err
}

func (p mockPersister) CreateJob(ctx context.Context, j core.Job) (core.Job, error) {
	return p.j, p.err
}

func (p mockPersister) DeleteJobsByWorkflowID(context.Context, string) error {
	return p.err
}

func (p mockPersister) GetJobs(ctx context.Context, pa core.PageArgs) (core.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, p.err
}

func (p mockPersister) GetJobsByWorkflowID(ctx context.Context, pa core.PageArgs, wid string) (core.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, p.err
}

func (p mockPersister) GetJobsByID(ctx context.Context, pa core.PageArgs, wid string, ids []string) (core.JobsPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.JobsPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.jp, p.err
}

func (p mockPersister) CreateVolume(context.Context, core.Volume) (core.Volume, error) {
	return p.v, p.err
}

func (p mockPersister) DeleteVolumesByWorkflowID(context.Context, string) error {
	return p.err
}

func (p mockPersister) GetVolumes(ctx context.Context, pa core.PageArgs) (core.VolumesPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.VolumesPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.vp, p.err
}

func (p mockPersister) GetVolumesByWorkflowID(ctx context.Context, pa core.PageArgs, wid string) (core.VolumesPage, error) {
	if got, want := pa, p.wantPA; !reflect.DeepEqual(got, want) {
		return core.VolumesPage{}, fmt.Errorf("got page args %v, want %v", got, want)
	}
	return p.vp, p.err
}

type mockIOFetcher struct {
	output string
	err    error
}

func (m mockIOFetcher) GetJobOutput(string) (string, error) {
	return m.output, m.err
}

type mockScheduler struct {
	err error
}

func (m mockScheduler) AddWorkflow(context.Context, core.Workflow, []core.Job, map[string]core.Volume) error {
	return m.err
}

type mockCore struct {
	p mockPersister
	f mockIOFetcher
	s mockScheduler
}

func getMockCore(mc mockCore) (*core.Core, error) {
	return core.New(&mc.p, &mc.f, &mc.s)
}
