// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package iomanager

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/rediskv"
)

// Config describes the IO manager configuration.
type Config struct {
	Version   string
	NATSConn  *nats.Conn
	RedisConn *rediskv.Connection
}

// IOManager contains the state of the IO Manager.
type IOManager struct {
	nc   *nats.Conn
	rc   *rediskv.Connection
	subs []*nats.Subscription
}

// New returns a new IO Manager.
func New(c Config) (m IOManager, err error) {
	return IOManager{
		c.NATSConn,
		c.RedisConn,
		nil,
	}, nil
}

// Start starts the IO Manager by initializing handlers for
// NATS subscriptions.
func (m IOManager) Start() {
	// Subscribe to relevant topics.
	if err := m.subscribe(); err != nil {
		logrus.WithError(err).Warn("failed to subscribe")
		return
	}
}

// Stop stops the IO Manager by putting NATS subscriptions
// in a draining state.
func (m IOManager) Stop() error {
	return m.unsubscribe()
}

// subscribe expresses interest in subjects that are relevant to the IO Manager.
func (m IOManager) subscribe() error {
	subs := []struct {
		subject string
		handler nats.MsgHandler
	}{
		{"job.*.output", m.jobOutputHandler},
	}
	for _, s := range subs {
		sub, err := m.nc.Subscribe(s.subject, s.handler)
		if err != nil {
			logrus.WithField("subject", s.subject).WithError(err).Warn("failed to subscribe")
			return err
		}
		logrus.WithField("subject", s.subject).Info("subscribed")

		m.subs = append(m.subs, sub)
	}
	return nil
}

// subscribe removes interest in subjects that are relevant to the IO Manager.
// NOTE: NATS will continue to handle callbacks until queue is empty.
func (m IOManager) unsubscribe() error {
	for _, s := range m.subs {
		err := s.Drain()
		if err != nil {
			logrus.WithField("subject", s.Subject).WithError(err).Warn("failed to unsubscribe")
			return err
		}
	}

	return nil
}

// NOTE: If multiple jobOutputHandlers are spun off, output for a job could be placed
// out of order in Redis.
func (m IOManager) jobOutputHandler(msg *nats.Msg) {
	// Parse subject for job ID
	s := strings.Split(msg.Subject, ".")
	if len(s) != 3 {
		logrus.Errorf("malformed job output subject: %s, skipping", msg.Subject)
	}

	id := s[1]
	if err := m.rc.Append(id, string(msg.Data)); err != nil {
		logrus.Errorf("failed to append job %s output: %v", id, err)
	}
}
