// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// getVersionHandler returns a version info handler.
func (s *Server) getVersionHandler() (http.Handler, error) {
	vr := struct {
		Version string `json:"version"`
	}{s.version}
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
func (*Server) getMetricsHandler() (http.Handler, error) {
	return promhttp.Handler(), nil
}
