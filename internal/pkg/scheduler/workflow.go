// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

const jobStartAckTimeout = time.Minute

// runJob runs a job to completion.
func (s *Scheduler) runJob(ctx context.Context, j model.Job) error {
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
		s.p.SetJobStatus(ctx, j.ID, msg.Status, msg.RC)
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

// runWorkflow runs a workflow to completion.
func (s *Scheduler) runWorkflow(w model.Workflow, jobs []model.Job) {
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

	for _, j := range jobs {
		ctx, cancel := context.WithTimeout(ctx, time.Minute) // TODO
		defer cancel()

		if err := s.runJob(ctx, j); err != nil {
			break
		}
	}

	s.p.SetWorkflowStatus(ctx, w.ID, "COMPLETED")
}

// AddWorkflow schedules a workflow for execution.
func (s *Scheduler) AddWorkflow(ctx context.Context, w model.Workflow, jobs []model.Job) error {
	s.p.SetWorkflowStatus(ctx, w.ID, "SCHEDULED")

	go s.runWorkflow(w, jobs)

	return nil
}
