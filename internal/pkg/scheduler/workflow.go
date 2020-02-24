// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/core"
)

const (
	jobStartAckTimeout = time.Minute
	volumeOpAckTimeout = time.Minute
)

// runJob runs a job to completion.
func (s *Scheduler) runJob(ctx context.Context, j core.Job) error {
	log := logrus.WithFields(logrus.Fields{
		"jobID":   j.ID,
		"jobName": j.Name,
	})
	log.Print("job starting")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("job completed")
	}(time.Now())

	jobFinished := make(chan struct{})

	// TODO: this should be a persistent subscription elsewhere.
	_, err := s.m.Subscribe(fmt.Sprintf("job.%v.finished", j.ID), func(msg struct {
		Status string
		RC     int
	}) {
		s.p.SetJobStatus(ctx, j.ID, msg.Status)
		s.p.SetJobExitCode(ctx, j.ID, msg.RC)
		close(jobFinished)
	})
	if err != nil {
		return err
	}

	var resp nats.Msg
	if err := s.m.Request("node.1.job.start", j, &resp, jobStartAckTimeout); err != nil {
		log.WithError(err).Print("failed to start job")
		return err
	}

	// Wait for response or timeout.
	select {
	case <-jobFinished:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// createVolume sets up a volume on an agent.
func (s *Scheduler) createVolume(ctx context.Context, v core.Volume) error {
	log := logrus.WithFields(logrus.Fields{
		"volumeID":   v.ID,
		"volumeName": v.Name,
		"volumeType": v.Type,
	})
	log.Print("creating volume")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("creation completed")
	}(time.Now())

	createFinished := make(chan error)

	// TODO: this should be a persistent subscription elsewhere.
	_, err := s.m.Subscribe(fmt.Sprintf("volume.%v.create", v.ID), func(msg struct {
		err error
	}) {
		createFinished <- msg.err
		close(createFinished)
	})
	if err != nil {
		return err
	}

	var resp nats.Msg
	if err := s.m.Request("node.1.volume.create", v, &resp, volumeOpAckTimeout); err != nil {
		log.WithError(err).Print("failed to create volume")
		return err
	}

	// Wait for response or timeout.
	select {
	case err := <-createFinished:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// deleteVolume tears down a volume on an agent.
func (s *Scheduler) deleteVolume(ctx context.Context, v core.Volume) error {
	log := logrus.WithFields(logrus.Fields{
		"volumeID":   v.ID,
		"volumeName": v.Name,
		"volumeType": v.Type,
	})
	log.Print("deleting volume")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("deletion completed")
	}(time.Now())

	deleteFinished := make(chan error)

	// TODO: this should be a persistent subscription elsewhere.
	_, err := s.m.Subscribe(fmt.Sprintf("volume.%v.delete", v.ID), func(msg struct {
		err error
	}) {
		deleteFinished <- msg.err
		close(deleteFinished)
	})
	if err != nil {
		return err
	}

	var resp nats.Msg
	if err := s.m.Request("node.1.volume.delete", v, &resp, volumeOpAckTimeout); err != nil {
		log.WithError(err).Print("failed to delete volume")
		return err
	}

	// Wait for response or timeout.
	select {
	case err := <-deleteFinished:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// runWorkflow runs a workflow to completion.
func (s *Scheduler) runWorkflow(w core.Workflow, jobs []core.Job, volumes map[string]core.Volume) {
	ctx := context.Background()

	log := logrus.WithFields(logrus.Fields{
		"workflowID":   w.ID,
		"workflowName": w.Name,
	})
	log.Print("workflow starting")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("workflow completed")
	}(time.Now())

	s.p.SetWorkflowStatus(ctx, w.ID, "RUNNING")

	// Bring up volumes on agent.
	for _, v := range volumes {
		ctx, cancel := context.WithTimeout(ctx, time.Minute) // TODO
		defer cancel()

		if err := s.createVolume(ctx, v); err != nil {
			break
		}
	}

	// Run jobs.
	for _, j := range jobs {
		ctx, cancel := context.WithTimeout(ctx, time.Minute) // TODO
		defer cancel()

		if err := s.runJob(ctx, j); err != nil {
			break
		}
	}

	// Tear down volumes on agent.
	for _, v := range volumes {
		ctx, cancel := context.WithTimeout(ctx, time.Minute) // TODO
		defer cancel()

		if err := s.deleteVolume(ctx, v); err != nil {
			break
		}
	}

	s.p.SetWorkflowStatus(ctx, w.ID, "COMPLETED")
}

// AddWorkflow schedules a workflow for execution.
func (s *Scheduler) AddWorkflow(ctx context.Context, w core.Workflow, jobs []core.Job, volumes map[string]core.Volume) error {
	s.p.SetWorkflowStatus(ctx, w.ID, "SCHEDULED")

	go s.runWorkflow(w, jobs, volumes)

	return nil
}
