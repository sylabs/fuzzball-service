// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sylabs/compute-service/internal/pkg/model"
	"gopkg.in/square/go-jose.v2"
)

func TestGetDiscoveryURIs(t *testing.T) {
	tests := []struct {
		name      string
		issuerURI string
		wantError bool
		wantURIs  []string
	}{
		{"BadURI", ":", true, nil},
		{"HostOnly", "https://example.com", false, []string{
			"https://example.com/.well-known/openid-configuration",
			"https://example.com/.well-known/oauth-authorization-server",
		}},
		{"EmptyPath", "https://example.com/", false, []string{
			"https://example.com/.well-known/openid-configuration",
			"https://example.com/.well-known/oauth-authorization-server",
		}},
		{"Path", "https://example.com/path", false, []string{
			"https://example.com/path/.well-known/openid-configuration",
			"https://example.com/.well-known/oauth-authorization-server/path",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDiscoveryURIs(tt.issuerURI)
			if (err != nil) != tt.wantError {
				t.Fatalf("got error %v, wantError %v", err, tt.wantError)
			}
			if !reflect.DeepEqual(got, tt.wantURIs) {
				t.Errorf("got %v, want %v", got, tt.wantURIs)
			}
		})
	}
}

type mockOAuthDisco struct {
	md model.AuthMetadata
}

var (
	testMetadata = model.AuthMetadata{
		Issuer:                                    "https://example.com/oauth2/default",
		AuthorizationEndpoint:                     "https://example.com/oauth2/default/v1/authorize",
		TokenEndpoint:                             "https://example.com/oauth2/default/v1/token",
		JWKSURI:                                   "https://example.com/oauth2/default/v1/keys",
		RegistrationEndpoint:                      "https://example.com/oauth2/v1/clients",
		ScopesSupported:                           []string{"openid", "profile", "email", "address", "phone", "offline_access"},
		ResponseTypesSupported:                    []string{"code", "token", "id_token", "code id_token", "code token", "id_token token", "code id_token token"},
		ResponseModesSupported:                    []string{"query", "fragment", "form_post", "okta_post_message"},
		GrantTypesSupported:                       []string{"authorization_code", "implicit", "refresh_token", "password", "client_credentials"},
		TokenEndpointAuthMethodsSupported:         []string{"client_secret_basic", "client_secret_post", "client_secret_jwt", "private_key_jwt", "none"},
		RevocationEndpoint:                        "https://example.com/oauth2/default/v1/revoke",
		RevocationEndpointAuthMethodsSupported:    []string{"client_secret_basic", "client_secret_post", "client_secret_jwt", "private_key_jwt", "none"},
		IntrospectionEndpoint:                     "https://example.com/oauth2/default/v1/introspect",
		IntrospectionEndpointAuthMethodsSupported: []string{"client_secret_basic", "client_secret_post", "client_secret_jwt", "private_key_jwt", "none"},
		CodeChallengeMethodsSupported:             []string{"S256"},
	}
)

func (m *mockOAuthDisco) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.md)
}

func TestGetAuthMetadata(t *testing.T) {
	expiredCtx, cancel := context.WithCancel(context.Background())
	cancel()

	ms := httptest.NewServer(&mockOAuthDisco{md: testMetadata})
	defer ms.Close()

	tests := []struct {
		name    string
		ctx     context.Context
		url     string
		wantErr bool
	}{
		{"ContextNil", nil, ms.URL, true},
		{"ContextExpired", expiredCtx, ms.URL, true},
		{"BadURL", context.Background(), "#", true},
		{"OK", context.Background(), ms.URL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAuthMetadata(tt.ctx, &http.Client{}, tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if want := testMetadata; !reflect.DeepEqual(got, want) {
					t.Errorf("got metadata %v, want %v", got, want)
				}
			}
		})
	}
}

type mockJWKS struct {
	keys jose.JSONWebKeySet
}

func (m *mockJWKS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.keys)
}

func TestGetKeySet(t *testing.T) {
	expiredCtx, cancel := context.WithCancel(context.Background())
	cancel()

	ms := httptest.NewServer(&mockJWKS{keys: testKeySet})
	defer ms.Close()

	tests := []struct {
		name    string
		ctx     context.Context
		url     string
		wantErr bool
	}{
		{"ContextNil", nil, ms.URL, true},
		{"ContextExpired", expiredCtx, ms.URL, true},
		{"BadURL", context.Background(), "#", true},
		{"OK", context.Background(), ms.URL, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getKeySet(tt.ctx, &http.Client{}, tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if want := testKeySet; !reflect.DeepEqual(got, want) {
					t.Errorf("got key set %v, want %v", got, want)
				}
			}
		})
	}
}
