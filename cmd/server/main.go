// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/app/server"
)

const (
	org  = "Sylabs"
	name = "Compute Server"
)

var (
	version = ""
)

var (
	httpAddr           = flag.String("http_addr", ":8080", "Address to bind HTTP")
	corsAllowedOrigins = flag.String("cors_allowed_origins", "*", "Comma-separated list of CORS allowed origins")
	corsDebug          = flag.Bool("cors_debug", false, "Enable CORS debugging")
	mongoURI           = flag.String("mongo_uri", "mongodb://localhost", "URI of MongoDB database")
	startupTime        = flag.Duration("startup_time", time.Minute, "Amount of time to wait for dependent services to become ready on startup")
)

// signalHandler catches SIGINT/SIGTERM to perform an orderly shutdown.
func signalHandler(s server.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	logrus.WithFields(logrus.Fields{
		"signal": (<-c).String(),
	}).Info("shutting down due to signal")

	if err := s.Stop(context.Background()); err != nil {
		logrus.WithError(err).Warning("shutdown failed")
	}
}

func main() {
	flag.Parse()

	log := logrus.WithFields(logrus.Fields{
		"org":  org,
		"name": name,
	})
	if version != "" {
		log = log.WithField("version", version)
	}
	log.Info("starting")

	c := server.Config{
		Version:            version,
		HTTPAddr:           *httpAddr,
		CORSAllowedOrigins: strings.Split(*corsAllowedOrigins, ","),
		CORSDebug:          *corsDebug,
		MongoURI:           *mongoURI,
	}

	ctx, cancel := context.WithTimeout(context.Background(), *startupTime)
	defer cancel()

	s, err := server.New(ctx, c)
	if err != nil {
		log.WithError(err).Fatal("failed to create server")
	}

	go signalHandler(s)

	s.Run()

	log.Info("stopping")
}
