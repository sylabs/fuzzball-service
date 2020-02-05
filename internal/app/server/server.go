// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/mongodb"
	"github.com/sylabs/compute-service/internal/pkg/resolver"
	"github.com/sylabs/compute-service/internal/pkg/scheduler"
	"github.com/sylabs/compute-service/internal/pkg/schema"
)

// Config describes server configuration.
type Config struct {
	Version            string
	HTTPAddr           string
	CORSAllowedOrigins []string
	CORSDebug          bool
	Persist            *mongodb.Connection
}

// Server contains the state of the server.
type Server struct {
	httpSrv *http.Server
	httpLn  net.Listener
	schema  *graphql.Schema
}

// New returns a new Server.
func New(ctx context.Context, c Config) (s Server, err error) {
	// Initialize scheduler.
	sched, err := scheduler.New(c.Persist)
	if err != nil {
		return Server{}, err
	}

	// Initialize GraphQL.
	sch, err := schema.String()
	if err != nil {
		return Server{}, fmt.Errorf("unable to init GraphQL schema: %w", err)
	}
	r, err := resolver.New(c.Persist, sched)
	if err != nil {
		return Server{}, err
	}
	schema, err := graphql.ParseSchema(sch, r)
	if err != nil {
		return Server{}, err
	}
	s.schema = schema

	// Set up HTTP server.
	h, err := s.NewRouter(c)
	if err != nil {
		return Server{}, err
	}
	s.httpSrv = &http.Server{
		Handler: loggingHandler(h),
	}

	// Start listening for HTTP.
	ln, err := net.Listen("tcp", c.HTTPAddr)
	if err != nil {
		return Server{}, err
	}
	defer func(ln net.Listener) {
		if err != nil {
			if err := ln.Close(); err != nil {
				logrus.WithError(err).Warning("error closing HTTP listener")
			}
		}
	}(ln)
	s.httpLn = ln
	logrus.WithField("addr", s.httpLn.Addr()).Info("listening")

	return s, nil
}

// Run is the main routine for the Server.
func (s Server) Run() {
	if err := s.httpSrv.Serve(s.httpLn); err != http.ErrServerClosed {
		logrus.WithError(err).Warning("serve error")
	}
}

// Stop is used to gracefully stop the Server.
func (s Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
