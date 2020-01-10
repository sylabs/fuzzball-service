// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"net"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/resolver"
	"github.com/sylabs/compute-service/internal/pkg/schema"
)

// Config describes server configuration.
type Config struct {
	Version            string
	HTTPAddr           string
	CORSAllowedOrigins []string
	CORSDebug          bool
	Persist            resolver.Persister
}

// Server contains the state of the server.
type Server struct {
	version            string
	corsAllowedOrigins []string
	corsDebug          bool
	httpSrv            *http.Server
	httpLn             net.Listener
	schema             *graphql.Schema
}

// New returns a new Server.
func New(ctx context.Context, c Config) (s Server, err error) {
	s = Server{
		version:            c.Version,
		corsAllowedOrigins: c.CORSAllowedOrigins,
		corsDebug:          c.CORSDebug,
	}

	// Initialize GraphQL.
	schema, err := graphql.ParseSchema(schema.String(), resolver.New(c.Persist))
	if err != nil {
		return Server{}, err
	}
	s.schema = schema

	// Set up HTTP server.
	h, err := NewRouter(&s)
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
