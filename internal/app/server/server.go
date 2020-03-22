// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/fuzzball-service/internal/pkg/core"
	"github.com/sylabs/fuzzball-service/internal/pkg/resolver"
	"github.com/sylabs/fuzzball-service/internal/pkg/schema"
	"gopkg.in/square/go-jose.v2"
)

// Config describes server configuration.
type Config struct {
	HTTPAddr                   string
	CORSAllowedOrigins         []string
	CORSDebug                  bool
	Core                       *core.Core
	OAuth2IssuerURI            string
	OAuth2Audience             string
	OAuth2Scopes               []string
	OAuth2PKCEClientID         string
	OAuth2PKCERedirectEndpoint string
}

// Server contains the state of the server.
type Server struct {
	httpSrv  *http.Server
	httpLn   net.Listener
	schema   *graphql.Schema
	authMeta core.AuthMetadata
	authKeys jose.JSONWebKeySet
}

// New returns a new Server.
func New(ctx context.Context, c Config) (s Server, err error) {
	hc := &http.Client{}

	// Discover OAuth 2.0 metadata.
	md, err := discoverAuthMetadata(ctx, hc, c.OAuth2IssuerURI)
	if err != nil {
		return Server{}, err
	}
	s.authMeta = md

	// Get OAuth key set.
	ks, err := getKeySet(ctx, hc, md.JWKSURI)
	if err != nil {
		return Server{}, err
	}
	s.authKeys = ks

	// Construct OAuth 2.0 configuration.
	oc := resolver.OAuth2Configuration{
		ClientCredentials: &resolver.ClientCredentialsConfig{
			TokenEndpoint: md.TokenEndpoint,
			Scopes:        c.OAuth2Scopes,
		},
	}
	if c.OAuth2PKCEClientID != "" && c.OAuth2PKCERedirectEndpoint != "" {
		oc.AuthCodePKCE = &resolver.AuthCodePKCEConfig{
			ClientID:              c.OAuth2PKCEClientID,
			AuthorizationEndpoint: md.AuthorizationEndpoint,
			TokenEndpoint:         md.TokenEndpoint,
			RedirectEndpoint:      c.OAuth2PKCERedirectEndpoint,
			Scopes:                c.OAuth2Scopes,
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"ClientID":         c.OAuth2PKCEClientID,
			"RedirectEndpoint": c.OAuth2PKCERedirectEndpoint,
		}).Warn("OAuth 2.0 Auth Code with PKCE disabled due to missing configuration value(s)")
	}

	// Initialize GraphQL.
	r, err := resolver.New(c.Core, oc)
	if err != nil {
		return Server{}, err
	}
	schema, err := schema.Get(r)
	if err != nil {
		return Server{}, fmt.Errorf("unable to init GraphQL schema: %w", err)
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
