// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/agent"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

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

	s.p.SetJobStatus(ctx, j.ID, "RUNNING", 0)

	rc, err := agent.RunJob(context.TODO(), j)
	if err != nil {
		s.p.SetJobStatus(ctx, j.ID, "FAILED", rc)
		return err
	}
	s.p.SetJobStatus(ctx, j.ID, "COMPLETED", rc)
	return nil
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
