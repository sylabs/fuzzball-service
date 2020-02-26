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

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/app/iomanager"
	"github.com/sylabs/compute-service/internal/app/server"
	"github.com/sylabs/compute-service/internal/pkg/mongodb"
	"github.com/sylabs/compute-service/internal/pkg/rediskv"
)

const (
	org  = "Sylabs"
	name = "Compute Server"

	dbName = "server"
)

var version = "unknown"

var (
	httpAddr           = flag.String("http-addr", ":8080", "Address to bind HTTP")
	corsAllowedOrigins = flag.String("cors-allowed-origins", "*", "Comma-separated list of CORS allowed origins")
	corsDebug          = flag.Bool("cors-debug", false, "Enable CORS debugging")
	mongoURI           = flag.String("mongo-uri", "mongodb://localhost", "URI of MongoDB database")
	startupTime        = flag.Duration("startup-time", time.Minute, "Amount of time to wait for dependent services to become ready on startup")
	natsURIs           = flag.String("nats-uris", "nats://localhost", "Comma-separated list of NATS server URIs")
	redisURI           = flag.String("redis-uri", "redis://localhost", "URI of Redis")
	oauth2IssuerURI    = flag.String("oauth2-issuer-uri", "https://dev-930666.okta.com/oauth2/default", "URI of OAuth 2.0 issuer")
	oauth2Audience     = flag.String("oauth2-audience", "api://default", "OAuth 2.0 audience expected in tokens")
)

// signalHandler catches SIGINT/SIGTERM to perform an orderly shutdown.
func signalHandler(nc *nats.Conn, s server.Server, m iomanager.IOManager) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	logrus.WithFields(logrus.Fields{
		"signal": (<-c).String(),
	}).Info("shutting down due to signal")

	if err := s.Stop(context.Background()); err != nil {
		logrus.WithError(err).Warning("server shutdown failed")
	}

	if err := m.Stop(); err != nil {
		logrus.WithError(err).Warning("IO manager shutdown failed")
	}

	// Drain nats connection before closing.
	if err := nc.Drain(); err != nil {
		logrus.WithError(err).Warning("starting nats connection draining failed")
	}

	// Wait for connection to drain and close.
	for nc.IsDraining() {
	}
}

// connectDB attempts to connect to the database.
func connectDB(ctx context.Context) (conn *mongodb.Connection, err error) {
	logrus.Info("connecting to database")
	defer func(t time.Time) {
		if err == nil {
			logrus.WithField("took", time.Since(t)).Info("database ready")
		}
	}(time.Now())

	return mongodb.NewConnection(ctx, *mongoURI, dbName)
}

// connectNATS attempts to connect to the NATS system.
func connectNATS(ctx context.Context) (nc *nats.Conn, err error) {
	logrus.Print("connecting to messaging system")
	defer func(t time.Time) {
		if err == nil {
			log := logrus.WithFields(logrus.Fields{
				"took":              time.Since(t),
				"connectedAddr":     nc.ConnectedAddr(),
				"connectedServerID": nc.ConnectedServerId(),
				"maxPayload":        nc.MaxPayload(),
			})
			if id, err := nc.GetClientID(); err == nil {
				// Log the client ID, if the server supports it.
				log = log.WithField("clientID", id)
			}
			log.Print("messaging system ready")
		}
	}(time.Now())

	o := nats.GetDefaultOptions()
	o.Servers = strings.Split(*natsURIs, ",")
	return o.Connect()
}

// connectRedis attempts to connect to redis.
func connectRedis() (*rediskv.Connection, error) {
	rc, err := rediskv.NewConnection(*redisURI)
	if err != nil {
		return nil, err
	}
	return rc, nil
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
	defer log.Info("stopped")

	// Context to control startup timeout.
	ctx, cancel := context.WithTimeout(context.Background(), *startupTime)
	defer cancel()

	// Connect to MongoDB.
	conn, err := connectDB(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to connect to database")
		return
	}
	defer func() {
		logrus.Info("disconnecting from database")
		if err := conn.Disconnect(context.Background()); err != nil {
			logrus.WithError(err).Warning("failed to disconnect from database")
		}
	}()

	// Connect to NATS.
	nc, err := connectNATS(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to connect to messaging system")
		return
	}
	defer func() {
		logrus.Info("disconnecting from messaging system")
		nc.Close()
	}()

	rc, err := connectRedis()
	if err != nil {
		logrus.WithError(err).Error("failed to connect to key value store")
		return
	}
	defer func() {
		logrus.Info("disconnecting from key value store")
		rc.Disconnect()
	}()

	// Spin up IO Manager.
	ioc := iomanager.Config{
		Version:   version,
		NATSConn:  nc,
		RedisConn: rc,
	}

	m, err := iomanager.New(ioc)
	if err != nil {
		logrus.WithError(err).Error("failed to create IO manager")
		return
	}
	m.Start()

	// Spin up server.
	c := server.Config{
		Version:            version,
		HTTPAddr:           *httpAddr,
		CORSAllowedOrigins: strings.Split(*corsAllowedOrigins, ","),
		CORSDebug:          *corsDebug,
		Persist:            conn,
		NATSConn:           nc,
		RedisConn:          rc,
		OAuth2IssuerURI:    *oauth2IssuerURI,
		OAuth2Audience:     *oauth2Audience,
	}
	s, err := server.New(ctx, c)
	if err != nil {
		logrus.WithError(err).Error("failed to create server")
		return
	}

	// Spin off signal handler to do graceful shutdown.
	go signalHandler(nc, s, m)

	// Main server routine.
	s.Run()
}
