// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// getVersion is used to retrieve version info from the server.
func (s *Server) getVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")

	vr := struct {
		Version string `json:"version"`
	}{s.version}
	if err := json.NewEncoder(w).Encode(vr); err != nil {
		logrus.WithError(err).Warning("failed to write response")
	}
}

// getMetrics is used to retrieve Prometheus metrics from the server.
func (*Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}
