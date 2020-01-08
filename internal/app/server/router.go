// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type handlerFunc func(*Server, http.ResponseWriter, *http.Request)

var routeConfigs = []struct {
	method  string
	pattern string
	handlerFunc
}{
	{http.MethodGet, "/version", (*Server).getVersion},
	{http.MethodGet, "/metrics", (*Server).getMetrics},
}

// handler returns an http.Handler that passes the Server as a receiver.
func (s *Server) handler(f handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(s, w, r)
	})
}

// NewRouter configures router and returns it.
func NewRouter(s *Server) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	for _, routeConfig := range routeConfigs {
		h := promhttp.InstrumentHandlerInFlight(httpRequestsInFlight,
			promhttp.InstrumentHandlerCounter(httpRequestsTotal,
				promhttp.InstrumentHandlerDuration(httpResponseTime,
					s.handler(routeConfig.handlerFunc),
				),
			),
		)

		router.
			Methods(routeConfig.method).
			Path(routeConfig.pattern).
			Handler(h)
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
	return c.Handler(router)
}
