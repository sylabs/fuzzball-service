// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type getHandlerFunc func(*Server) (http.Handler, error)

var routeConfigs = []struct {
	pattern string
	getHandlerFunc
}{
	{"/version", (*Server).getVersionHandler},
	{"/metrics", (*Server).getMetricsHandler},
}

// NewRouter configures router and returns it.
func NewRouter(s *Server) (http.Handler, error) {
	mux := http.NewServeMux()

	for _, routeConfig := range routeConfigs {
		// Get handler.
		h, err := routeConfig.getHandlerFunc(s)
		if err != nil {
			return nil, err
		}

		// Instrument with Prometheus.
		h = promhttp.InstrumentHandlerInFlight(httpRequestsInFlight,
			promhttp.InstrumentHandlerCounter(httpRequestsTotal,
				promhttp.InstrumentHandlerDuration(httpResponseTime, h),
			),
		)

		// Add to mux.
		mux.Handle(routeConfig.pattern, h)
	}

	// Implement CORS specification for all routes.
	c := cors.New(cors.Options{
		AllowedOrigins: s.corsAllowedOrigins,
		AllowedMethods: []string{
			http.MethodDelete,
			http.MethodGet,
			http.MethodOptions,
			http.MethodPost,
		},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Language",
			"Authorization",
			"Content-Language",
			"Content-Type",
		},
		AllowCredentials: true,
		Debug:            s.corsDebug,
	})
	return c.Handler(mux), nil
}
