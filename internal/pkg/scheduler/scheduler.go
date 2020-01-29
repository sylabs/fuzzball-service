// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/agent"
	"github.com/sylabs/compute-service/internal/pkg/model"
)

// Persister is the interface that describes what is needed to persist scheduler data.
type Persister interface {
	SetWorkflowStatus(context.Context, string, string) error
	SetJobStatus(context.Context, string, string, int) error
}

func runJob(ctx context.Context, p Persister, j model.Job) error {
	log := logrus.WithFields(logrus.Fields{
		"jobID":   j.ID,
		"jobName": j.Name,
	})
	log.Print("job starting")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("job completed")
	}(time.Now())

	p.SetJobStatus(ctx, j.ID, "RUNNING", 0)

	rc, err := agent.RunJob(context.TODO(), j)
	if err != nil {
		p.SetJobStatus(ctx, j.ID, "FAILED", rc)
		return err
	}
	p.SetJobStatus(ctx, j.ID, "COMPLETED", rc)
	return nil
}

// runWorkflow runs a workflow to completion.
func runWorkflow(p Persister, w model.Workflow, jobs []model.Job) {
	ctx := context.Background()

	log := logrus.WithFields(logrus.Fields{
		"workflowID":   w.ID,
		"workflowName": w.Name,
	})
	log.Print("workflow starting")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("workflow completed")
	}(time.Now())

	p.SetWorkflowStatus(ctx, w.ID, "RUNNING")

	for _, j := range jobs {
		if err := runJob(ctx, p, j); err != nil {
			break
		}
	}

	p.SetWorkflowStatus(ctx, w.ID, "COMPLETED")
}

// AddWorkflow schedules a workflow for execution.
func AddWorkflow(ctx context.Context, p Persister, w model.Workflow, jobs []model.Job) error {
	p.SetWorkflowStatus(ctx, w.ID, "SCHEDULED")

	go runWorkflow(p, w, jobs)

	return nil
}
