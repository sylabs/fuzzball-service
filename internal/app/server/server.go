// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "server"

// Config describes server configuration.
type Config struct {
	Version            string
	HTTPAddr           string
	CORSAllowedOrigins []string
	CORSDebug          bool
	MongoURI           string
}

// Server contains the state of the server.
type Server struct {
	version            string
	corsAllowedOrigins []string
	corsDebug          bool
	httpSrv            *http.Server
	httpLn             net.Listener
	db                 *mongo.Database
}

// newDB connects to the DB, and returns a Database.
func newDB(ctx context.Context, mongoURI string) (db *mongo.Database, err error) {
	logrus.Info("connecting to database")
	defer func(t time.Time) {
		if err == nil {
			logrus.WithField("took", time.Since(t)).Info("database ready")
		}
	}(time.Now())

	o := options.Client().ApplyURI(mongoURI)
	if err := o.Validate(); err != nil {
		return nil, err
	}
	mc, err := mongo.NewClient(o)
	if err != nil {
		return nil, err
	}
	if err := mc.Connect(ctx); err != nil {
		return nil, err
	}
	if err := mc.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return mc.Database(dbName), nil
}

// New returns a new Server.
func New(ctx context.Context, c Config) (s Server, err error) {
	s = Server{
		version:            c.Version,
		corsAllowedOrigins: c.CORSAllowedOrigins,
		corsDebug:          c.CORSDebug,
		httpSrv: &http.Server{
			Handler: loggingHandler(NewRouter(&s)),
		},
	}

	// Connect to database.
	db, err := newDB(ctx, c.MongoURI)
	if err != nil {
		return Server{}, err
	}
	defer func(db *mongo.Database) {
		if err != nil {
			if err := db.Client().Disconnect(ctx); err != nil {
				logrus.WithError(err).Warning("disconnect error")
			}
		}
	}(db)
	s.db = db

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
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Warning("shutdown error")
	}

	if err := s.db.Client().Disconnect(ctx); err != nil {
		logrus.WithError(err).Warning("disconnect error")
	}
	return nil
}
