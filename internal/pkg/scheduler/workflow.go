// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package scheduler

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/fuzzball-service/internal/pkg/core"
	scs "github.com/sylabs/scs-library-client/client"
)

const (
	jobStartAckTimeout        = time.Minute
	volumeOpAckTimeout        = time.Minute
	cacheOpAckTimeout         = time.Minute
	imageDownloadOpAckTimeout = 10 * time.Minute
)

// agentCacheInfo describes the state of the cache on the agent
type agentCacheInfo struct {
	Cached bool
	Hash   string
}

// runJob runs a job to completion.
func (s *Scheduler) runJob(ctx context.Context, j core.Job, ac agentCacheInfo) error {
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

	jobInfo := struct {
		core.Job
		agentCacheInfo
	}{
		j,
		ac,
	}

	var resp nats.Msg
	if err := s.m.Request("node.1.job.start", jobInfo, &resp, jobStartAckTimeout); err != nil {
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

type image struct {
	URI string
}

// imageDownload pull an image to the cache on the agent.
func (s *Scheduler) imageDownload(ctx context.Context, i image) error {
	log := logrus.WithFields(logrus.Fields{
		"imageURI": i.URI,
	})
	log.Print("downloading image")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("download completed")
	}(time.Now())

	downloadFinished := make(chan error)

	// TODO: this should be a persistent subscription elsewhere.
	sub, err := s.m.Subscribe("image.download", func(msg struct {
		err error
	}) {
		downloadFinished <- msg.err
		close(downloadFinished)
	})
	if err != nil {
		return err
	}
	sub.AutoUnsubscribe(1)

	var resp nats.Msg
	if err := s.m.Request("node.1.image.download", i, &resp, imageDownloadOpAckTimeout); err != nil {
		log.WithError(err).Print("failed to download image")
		return err
	}

	// Wait for response or timeout.
	select {
	case err := <-downloadFinished:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// imageCached checks the agent image cache for existance of an image based on its hash.
func (s *Scheduler) imageCached(ctx context.Context, hash string) (bool, error) {
	log := logrus.WithFields(logrus.Fields{
		"hash": hash,
	})
	log.Print("checking agent image cache")
	defer func(t time.Time) {
		log.WithField("took", time.Since(t)).Print("agent image cache check completed")
	}(time.Now())

	checkFinished := make(chan bool)

	// TODO: this should be a persistent subscription elsewhere.
	sub, err := s.m.Subscribe("image.cached", func(msg struct {
		exists bool
	}) {
		checkFinished <- msg.exists
		close(checkFinished)
	})
	if err != nil {
		return false, err
	}
	sub.AutoUnsubscribe(1)

	var resp nats.Msg
	if err := s.m.Request("node.1.image.cached", hash, &resp, cacheOpAckTimeout); err != nil {
		log.WithError(err).Print("failed to get cache data")
		return false, err
	}

	// Wait for response or timeout.
	select {
	case exists := <-checkFinished:
		return exists, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (s *Scheduler) prepAgent(ctx context.Context, j core.Job) (ac agentCacheInfo, err error) {
	logrus.Print("preparing agent")
	defer func(t time.Time) {
		logrus.WithField("took", time.Since(t)).Print("agent prepared")
	}(time.Now())
	// skip caching for sources other than the library
	if !strings.HasPrefix(j.Image, "library:") {
		return ac, nil
	}

	// Parse image ref
	r, err := scs.Parse(j.Image)
	if err != nil {
		return ac, err
	}

	if len(r.Tags) == 0 {
		r.Tags = []string{"latest"}
	}

	// Point library client to specific library if included in uri
	var scsConf *scs.Config
	if r.Host != "" {
		scsConf = &scs.Config{BaseURL: "https://" + r.Host}
	}

	// Initialize library client
	client, err := scs.NewClient(scsConf)
	if err != nil {
		logrus.WithError(err).Warnf("could not initialize library client")
		return ac, err
	}

	// Get library image metadata using the path and tag of the uri
	imageRef := r.Path + ":" + r.Tags[0]
	// TODO: this should use agent GOARCH, should be part of agent check-in
	meta, err := client.GetImage(ctx, runtime.GOARCH, imageRef)
	if err != nil {
		logrus.WithError(err).Warnf("could not fetch image metadata")
		return ac, err
	}

	// Check image in cache by hash
	cached, err := s.imageCached(ctx, meta.Hash)
	if err != nil {
		logrus.WithError(err).Warnf("while checking agent cache")
		return ac, err
	}

	// Set reference tag to be image hash
	r.Tags = []string{meta.Hash}
	if !cached {
		// Have agent download image by hash
		err := s.imageDownload(ctx, image{r.String()})
		if err != nil {
			logrus.WithError(err).Warnf("while downloading image to agent cache")
			return ac, err
		}
	}

	ac.Cached = true
	ac.Hash = meta.Hash

	return ac, nil
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

		ci, err := s.prepAgent(ctx, j)
		if err != nil {
			break
		}

		// NOTE: default to singularity image pulling for non-library images for now
		if err := s.runJob(ctx, j, ci); err != nil {
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
