// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type getHandlerFunc func(*Server, Config) (http.Handler, error)

var routeConfigs = []struct {
	pattern string
	getHandlerFunc
}{
	{"/version", (*Server).getVersionHandler},
	{"/metrics", (*Server).getMetricsHandler},
	{"/graphql", (*Server).getGraphQLHandler},
	{"/graphiql", (*Server).getGraphiQLHandler},
}

// NewRouter configures router and returns it.
func (s *Server) NewRouter(c Config) (http.Handler, error) {
	mux := http.NewServeMux()

	for _, routeConfig := range routeConfigs {
		// Get handler.
		h, err := routeConfig.getHandlerFunc(s, c)
		if err != nil {
			return nil, err
		}

		// Add JWT middleware.
		h = s.tokenHandler(c, h)

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
	cors := cors.New(cors.Options{
		AllowedOrigins: c.CORSAllowedOrigins,
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
		Debug:            c.CORSDebug,
	})
	return cors.Handler(mux), nil
}
