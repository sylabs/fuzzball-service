// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
