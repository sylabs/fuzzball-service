// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package token

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		wantCode int
	}{
		{"OK", testToken, http.StatusOK},
		{"NoClaims", testTokenNoClaims, http.StatusOK},
		{"InvalidAudience", testTokenInvalidAudience, http.StatusUnauthorized},
		{"InvalidIssuer", testTokenInvalidIssuer, http.StatusUnauthorized},
		{"Expired", testTokenExpired, http.StatusUnauthorized},
		{"NoAuthHeader", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiddleware(MiddlewareOptions{
				Audience: testClaims.Audience,
				Issuer:   testClaims.Issuer,
				KeyFunc: func(t *jwt.Token) (interface{}, error) {
					return testSigningKey, nil
				},
			})
			h := m.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := FromContext(r.Context()); !ok {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}))

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, r)

			if got := rr.Code; got != tt.wantCode {
				t.Errorf("got code %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		want       string
		wantErr    bool
	}{
		{"NoToken", "Bearer", "", true},
		{"NoAuthHeader", "", "", true},
		{"OK", "Bearer TOKEN", "TOKEN", false},
		{"Whitespace", " Bearer  TOKEN ", "TOKEN", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				r.Header.Set("Authorization", tt.authHeader)
			}

			got, err := getBearerToken(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got token %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyJWT(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"OK", testToken, false},
		{"NoClaims", testTokenNoClaims, false},
		{"InvalidAudience", testTokenInvalidAudience, true},
		{"InvalidIssuer", testTokenInvalidIssuer, true},
		{"Expired", testTokenExpired, true},
		{"NoAuthHeader", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiddleware(MiddlewareOptions{
				Audience: testClaims.Audience,
				Issuer:   testClaims.Issuer,
				KeyFunc: func(t *jwt.Token) (interface{}, error) {
					return testSigningKey, nil
				},
			})

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.token != "" {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}

			if err := m.verifyJWT(r); (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
