// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	exampleURL = "http://example.com"
	otherURL   = "http://other.net"
)

// testRouteConfigs is a list of routes to test, along with the expected HTTP status code that will
// be returned when an empty request is sent with no authorization header.
var testRouteConfigs = []struct {
	name         string
	method       string
	path         string
	expectedCode int
}{
	{"GetMetrics", http.MethodGet, "/metrics", http.StatusOK},
	{"PostGraphQL", http.MethodPost, "/graphql", http.StatusBadRequest},
	{"GetGraphiQL", http.MethodGet, "/graphiql", http.StatusOK},
}

func TestRouteConfigs(t *testing.T) {
	if have, want := len(testRouteConfigs), len(routeConfigs); have != want {
		t.Errorf("have %v route configs, want %v", have, want)
	}
}

func TestRouter(t *testing.T) {
	sr := Server{}
	cfg := Config{}
	h, err := sr.NewRouter(cfg)
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	s := httptest.NewServer(h)
	defer s.Close()

	c := http.Client{}
	for _, rtt := range testRouteConfigs {
		t.Run(rtt.name, func(t *testing.T) {
			req, err := http.NewRequest(rtt.method, fmt.Sprintf("%s%s", s.URL, rtt.path), nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			res, err := c.Do(req)
			if err != nil {
				t.Fatalf("failed to do request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != rtt.expectedCode {
				t.Fatalf("unexpected status code: %v/%v", res.StatusCode, rtt.expectedCode)
			}
		})
	}
}

func TestRouterCORS(t *testing.T) {
	c := http.Client{}

	tests := []struct {
		name           string
		allowedOrigins []string
		setOrigin      bool
		expectCORS     bool
		expectedOrigin string
	}{
		{"NoOriginWildCard", []string{"*"}, false, false, ""},
		{"NoOriginSingle", []string{exampleURL}, false, false, ""},
		{"NoOriginMultiple", []string{exampleURL, otherURL}, false, false, ""},
		{"MatchWildcard", []string{"*"}, true, true, "*"},
		{"MatchSingle", []string{exampleURL}, true, true, exampleURL},
		{"MatchMultiple", []string{exampleURL, otherURL}, true, true, exampleURL},
		{"NoMatch", []string{otherURL}, true, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := Server{}
			cfg := Config{
				CORSAllowedOrigins: tt.allowedOrigins,
			}
			h, err := sr.NewRouter(cfg)
			if err != nil {
				t.Fatalf("failed to create router: %v", err)
			}

			s := httptest.NewServer(h)
			defer s.Close()

			for _, rtt := range testRouteConfigs {
				t.Run(rtt.name, func(t *testing.T) {
					req, err := http.NewRequest(rtt.method, fmt.Sprintf("%s%s", s.URL, rtt.path), nil)
					if err != nil {
						t.Fatalf("failed to create request: %v", err)
					}
					if tt.setOrigin {
						req.Header.Set("Origin", exampleURL)
					}

					res, err := c.Do(req)
					if err != nil {
						t.Fatalf("failed to do request: %v", err)
					}
					defer res.Body.Close()

					if res.StatusCode != rtt.expectedCode {
						t.Fatalf("unexpected status code: %v/%v", res.StatusCode, rtt.expectedCode)
					}

					// The "Vary" header should always contain "Origin".
					if got, want := res.Header.Get("Vary"), "Origin"; got != want {
						t.Errorf("got Vary header '%v', want '%v'", got, want)
					}

					// If CORS headers are expected, verify they contain expected values.
					hdrs := map[string]string{
						"Access-Control-Allow-Origin":      tt.expectedOrigin,
						"Access-Control-Allow-Credentials": "true",
					}
					for k, v := range hdrs {
						if _, ok := res.Header[k]; ok != tt.expectCORS {
							t.Errorf("%v present is %v, want %v", k, ok, tt.expectCORS)
						} else if ok && tt.expectCORS {
							if got := res.Header.Get(k); got != v {
								t.Errorf("%v value is %v, want %v", k, got, v)
							}
						}
					}
				})
			}
		})
	}
}

func TestRouterCORSPreflight(t *testing.T) {
	c := http.Client{}

	tests := []struct {
		name             string
		allowedOrigins   []string
		setOrigin        bool
		setRequestMethod bool
		requestHeaders   string
		expectCORS       bool
		expectedOrigin   string
	}{
		{"NoOriginWildCard", []string{"*"}, false, true, "", false, ""},
		{"NoOriginSingle", []string{exampleURL}, false, true, "", false, ""},
		{"NoOriginMultiple", []string{exampleURL, otherURL}, false, true, "", false, ""},
		{"MatchWildcard", []string{"*"}, true, true, "", true, "*"},
		{"MatchSingle", []string{exampleURL}, true, true, "", true, exampleURL},
		{"MatchMultiple", []string{exampleURL, otherURL}, true, true, "", true, exampleURL},
		{"MatchWildcardHeaders", []string{"*"}, true, true, "Authorization", true, "*"},
		{"MatchSingleHeaders", []string{exampleURL}, true, true, "Authorization", true, exampleURL},
		{"MatchMultipleHeaders", []string{exampleURL, otherURL}, true, true, "Authorization", true, exampleURL},
		{"NoMatch", []string{otherURL}, true, true, "", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := Server{}
			cfg := Config{
				CORSAllowedOrigins: tt.allowedOrigins,
			}
			h, err := sr.NewRouter(cfg)
			if err != nil {
				t.Fatalf("failed to create router: %v", err)
			}

			s := httptest.NewServer(h)
			defer s.Close()

			for _, rtt := range testRouteConfigs {
				t.Run(rtt.name, func(t *testing.T) {
					req, err := http.NewRequest(http.MethodOptions, fmt.Sprintf("%s%s", s.URL, rtt.path), nil)
					if err != nil {
						t.Fatalf("failed to create request: %v", err)
					}
					if tt.setOrigin {
						req.Header.Set("Origin", exampleURL)
					}
					if tt.setRequestMethod {
						req.Header.Set("Access-Control-Request-Method", rtt.method)
					}
					if tt.requestHeaders != "" {
						req.Header.Set("Access-Control-Request-Headers", tt.requestHeaders)
					}

					res, err := c.Do(req)
					if err != nil {
						t.Fatalf("failed to do request: %v", err)
					}
					defer res.Body.Close()

					// Preflight check should return 200 OK
					if got, want := res.StatusCode, http.StatusOK; got != want {
						t.Fatalf("unexpected status code: %v/%v", got, want)
					}

					// The "Vary" header should always contain "Origin".
					if got, want := res.Header.Get("Vary"), "Origin"; got != want {
						t.Errorf("got Vary header '%v', want '%v'", got, want)
					}

					// If CORS headers are expected, verify they contain expected values.
					hdrs := map[string]string{
						"Access-Control-Allow-Origin":      tt.expectedOrigin,
						"Access-Control-Allow-Methods":     rtt.method,
						"Access-Control-Allow-Credentials": "true",
					}
					if tt.requestHeaders != "" {
						hdrs["Access-Control-Allow-Headers"] = tt.requestHeaders
					}
					for k, v := range hdrs {
						if _, ok := res.Header[k]; ok != tt.expectCORS {
							t.Errorf("%v present is %v, want %v", k, ok, tt.expectCORS)
						} else if ok && tt.expectCORS {
							if got := res.Header.Get(k); got != v {
								t.Errorf("%v value is %v, want %v", k, got, v)
							}
						}
					}
				})
			}
		})
	}
}

func TestRouterNotFound(t *testing.T) {
	sr := Server{}
	h, err := sr.NewRouter(Config{})
	if err != nil {
		t.Fatalf("failed to create router: %v", err)
	}

	s := httptest.NewServer(h)
	defer s.Close()

	res, err := http.Get(fmt.Sprintf("%s%s", s.URL, "/not/a/valid/path"))
	if err != nil {
		t.Fatalf("failed to do request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status code: %v/%v", res.StatusCode, http.StatusNotFound)
	}
}
