// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricsNamespace = "fuzzballserver"
)

var (
	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "http_requests_in_flight",
		Help:      "Current number of HTTP requests being served.",
	})
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricsNamespace,
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests by status code and method.",
	}, []string{"code", "method"})
	httpResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricsNamespace,
		Name:      "http_response_time_seconds",
		Help:      "Histogram of HTTP response time in seconds.",
	}, []string{})
)
