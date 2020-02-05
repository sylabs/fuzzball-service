// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/friendsofgo/graphiql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// getVersionHandler returns a version info handler.
func (s *Server) getVersionHandler(c Config) (http.Handler, error) {
	vr := struct {
		Version string `json:"version"`
	}{c.Version}
	b, err := json.Marshal(vr)
	if err != nil {
		return nil, err
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			logrus.WithError(err).Warning("failed to write response")
		}
	}
	return http.HandlerFunc(h), nil
}

// getMetricsHandler returns a Prometheus metrics handler.
func (*Server) getMetricsHandler(c Config) (http.Handler, error) {
	ph := promhttp.Handler()

	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ph.ServeHTTP(w, r)
	}
	return http.HandlerFunc(h), nil
}

// getGraphQLHandler returns a GraphQL handler.
func (s *Server) getGraphQLHandler(c Config) (http.Handler, error) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res := s.schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			logrus.WithError(err).Warning("failed to write response")
		}
	}
	return http.HandlerFunc(h), nil
}

// getGraphiQLHandler returns a GraphiQL handler.
func (s *Server) getGraphiQLHandler(c Config) (http.Handler, error) {
	return graphiql.NewGraphiqlHandler("/graphql")
}
