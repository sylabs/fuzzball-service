// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		version  string
		wantCode int
	}{
		{"GetVersion", http.MethodGet, "1.0.0", http.StatusOK},
		{"PostVersion", http.MethodPost, "", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Server{}
			cfg := Config{
				Version: tt.version,
			}
			h, err := s.getVersionHandler(cfg)
			if err != nil {
				t.Fatalf("failed to get handler: %v", err)
			}

			rr := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "/version", nil)

			h.ServeHTTP(rr, r)

			if got, want := rr.Code, tt.wantCode; got != want {
				t.Fatalf("got code %v, want %v", got, want)
			}

			if rr.Code == http.StatusOK {
				var vr struct {
					Version string `json:"version"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&vr); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				if got, want := vr.Version, tt.version; got != want {
					t.Errorf("got version %v, want %v", got, want)
				}
			}
		})
	}
}

func TestGetMetrics(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		wantCode int
	}{
		{"GetMetrics", http.MethodGet, http.StatusOK},
		{"PostMetrics", http.MethodPost, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{}

			h, err := s.getMetricsHandler(Config{})
			if err != nil {
				t.Fatalf("failed to get handler: %v", err)
			}

			rr := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, "/metrics", nil)

			h.ServeHTTP(rr, r)

			if got, want := rr.Code, tt.wantCode; got != want {
				t.Fatalf("got code %v, want %v", got, want)
			}
		})
	}
}
