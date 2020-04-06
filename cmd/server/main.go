// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/blang/semver"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/sylabs/fuzzball-service/internal/app/iomanager"
	"github.com/sylabs/fuzzball-service/internal/app/server"
	"github.com/sylabs/fuzzball-service/internal/pkg/core"
	"github.com/sylabs/fuzzball-service/internal/pkg/mongodb"
	"github.com/sylabs/fuzzball-service/internal/pkg/rediskv"
	"github.com/sylabs/fuzzball-service/internal/pkg/scheduler"
)

const (
	org  = "Sylabs"
	name = "Fuzzball Server"

	dbName = "server"

	keyStartupTime                = "startup-time"
	keyHTTPAddr                   = "http-addr"
	keyCORSAllowedOrigins         = "cors-allowed-origins"
	keyCORSDebug                  = "cors-debug"
	keyMongoURI                   = "mongo-uri"
	keyNatsURIs                   = "nats-uris"
	keyRedisURI                   = "redis-uri"
	keyOAuth2IssuerURI            = "oauth2-issuer-uri"
	keyOAuth2Audience             = "oauth2-audience"
	keyOAuth2Scopes               = "oauth2-scopes"
	keyOAuth2PKCEClientID         = "oauth2-pkce-client-id"
	keyOAuth2PKCERedirectEndpoint = "oauth2-pkce-redirect-endpoint"
)

// Values set during build.
var (
	builtAt      = ""
	gitCommit    = ""
	gitTreeState = ""
	gitVersion   = ""
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
func connectDB(ctx context.Context, uri string) (mc *mongodb.Connection, err error) {
	logrus.Info("connecting to database")
	defer func(t time.Time) {
		if err == nil {
			logrus.WithField("took", time.Since(t)).Info("database ready")
		}
	}(time.Now())

	return mongodb.NewConnection(ctx, uri, dbName)
}

// connectNATS attempts to connect to the NATS system.
func connectNATS(ctx context.Context, uris []string) (nc *nats.Conn, err error) {
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
	o.Servers = uris
	return o.Connect()
}

// connectRedis attempts to connect to redis.
func connectRedis(uri string) (*rediskv.Connection, error) {
	rc, err := rediskv.NewConnection(uri)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// getFlagSet declares and parses the command line flags.
func getFlagSet() *pflag.FlagSet {
	fs := pflag.CommandLine
	fs.Duration(keyStartupTime, time.Minute, "Amount of time to wait for dependent services to become ready on startup")
	fs.String(keyHTTPAddr, ":8080", "Address to bind HTTP")
	fs.StringSlice(keyCORSAllowedOrigins, []string{"*"}, "Comma-separated list of CORS allowed origins")
	fs.Bool(keyCORSDebug, false, "Enable CORS debugging")
	fs.String(keyMongoURI, "mongodb://localhost", "URI of MongoDB database")
	fs.StringSlice(keyNatsURIs, []string{"nats://localhost"}, "Comma-separated list of NATS server URIs")
	fs.String(keyRedisURI, "redis://localhost", "URI of Redis")
	fs.String(keyOAuth2IssuerURI, "https://dev-930666.okta.com/oauth2/default", "URI of OAuth 2.0 issuer")
	fs.String(keyOAuth2Audience, "api://default", "OAuth 2.0 audience expected in tokens")
	fs.StringSlice(keyOAuth2Scopes, []string{"openid", "offline_access"}, "Recommended scope(s) for OAuth 2.0 clients to request")
	fs.String(keyOAuth2PKCEClientID, "", "Client ID for OAuth 2.0 clients to use for Authorization Code flow with PKCE")
	fs.String(keyOAuth2PKCERedirectEndpoint, "http://localhost:9876/authorization/callback", "Callback URL for OAuth 2.0 clients to use for Authorization Code flow with PKCE")

	fs.Parse(os.Args[1:])

	return fs
}

// getConfig gets a Viper instance to retrieve configuration.
func getConfig() (*viper.Viper, error) {
	v := viper.New()

	// Bind command line flags.
	if err := v.BindPFlags(getFlagSet()); err != nil {
		return nil, err
	}

	// Set up to use environment.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return v, nil
}

// getCore returns an initilized Core.
func getCore(mc *mongodb.Connection, nc *nats.Conn, rc *rediskv.Connection) (*core.Core, error) {
	// Encoded NATS connection.
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	// Initialize scheduler.
	sched, err := scheduler.New(ec, mc, rc)
	if err != nil {
		return nil, err
	}

	// Build up core options.
	opts := [](func(*core.Core) error){}
	if t, err := time.Parse(time.RFC3339, builtAt); err == nil {
		opts = append(opts, core.OptBuiltAt(t))
	}
	if gitCommit != "" {
		opts = append(opts, core.OptGitCommit(gitCommit))
	}
	if gitTreeState != "" {
		opts = append(opts, core.OptGitTreeState(gitTreeState))
	}
	if v, err := semver.Parse(gitVersion); err == nil {
		opts = append(opts, core.OptGitVersion(v))
	}

	// Initialize core.
	return core.New(mc, rc, sched, opts...)
}

func main() {
	log := logrus.WithFields(logrus.Fields{
		"org":  org,
		"name": name,
	})
	if builtAt != "" {
		log = log.WithField("builtAt", builtAt)
	}
	if gitCommit != "" {
		log = log.WithField("gitCommit", gitCommit)
	}
	if gitTreeState != "" {
		log = log.WithField("gitTreeState", gitTreeState)
	}
	if gitVersion != "" {
		log = log.WithField("gitVersion", gitVersion)
	}
	log.Info("starting")
	defer log.Info("stopped")

	// Create viper instance, which holds configuration.
	cfg, err := getConfig()
	if err != nil {
		logrus.WithError(err).Error("failed to get configuration")
		return
	}

	// Context to control startup timeout.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GetDuration(keyStartupTime))
	defer cancel()

	// Connect to MongoDB.
	mc, err := connectDB(ctx, cfg.GetString(keyMongoURI))
	if err != nil {
		logrus.WithError(err).Error("failed to connect to database")
		return
	}
	defer func() {
		logrus.Info("disconnecting from database")
		if err := mc.Disconnect(context.Background()); err != nil {
			logrus.WithError(err).Warning("failed to disconnect from database")
		}
	}()

	// Connect to NATS.
	nc, err := connectNATS(ctx, cfg.GetStringSlice(keyNatsURIs))
	if err != nil {
		logrus.WithError(err).Error("failed to connect to messaging system")
		return
	}
	defer func() {
		logrus.Info("disconnecting from messaging system")
		nc.Close()
	}()

	// Connect to Redis.
	rc, err := connectRedis(cfg.GetString(keyRedisURI))
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
		NATSConn:  nc,
		RedisConn: rc,
	}
	m, err := iomanager.New(ioc)
	if err != nil {
		logrus.WithError(err).Error("failed to create IO manager")
		return
	}
	m.Start()

	// Get core.
	c, err := getCore(mc, nc, rc)
	if err != nil {
		logrus.WithError(err).Error("failed to get core")
		return
	}

	// Set up server configuration.
	sc := server.Config{
		HTTPAddr:                   cfg.GetString(keyHTTPAddr),
		CORSAllowedOrigins:         cfg.GetStringSlice(keyCORSAllowedOrigins),
		CORSDebug:                  cfg.GetBool(keyCORSDebug),
		OAuth2IssuerURI:            cfg.GetString(keyOAuth2IssuerURI),
		OAuth2Audience:             cfg.GetString(keyOAuth2Audience),
		OAuth2Scopes:               cfg.GetStringSlice(keyOAuth2Scopes),
		OAuth2PKCEClientID:         cfg.GetString(keyOAuth2PKCEClientID),
		OAuth2PKCERedirectEndpoint: cfg.GetString(keyOAuth2PKCERedirectEndpoint),
	}

	// Spin up server.
	s, err := server.New(ctx, c, sc)
	if err != nil {
		logrus.WithError(err).Error("failed to create server")
		return
	}

	// Spin off signal handler to do graceful shutdown.
	go signalHandler(nc, s, m)

	// Main server routine.
	s.Run()
}
