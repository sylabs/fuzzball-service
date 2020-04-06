// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package resolver

// OAuth2Configuration represents OAuth 2.0 configuration.
type OAuth2Configuration struct {
	AuthCodePKCE      *AuthCodePKCEConfig      // Configuration for Authorization Code Flow with Proof Key for Code Exchange (if supported).
	ClientCredentials *ClientCredentialsConfig // Configuration for Client Credentials Flow (if supported).
}

// AuthCodePKCEConfig contains OAuth 2.0 configuration for the Authorization Code Flow with Proof
// Key for Code Exchange (PKCE).
type AuthCodePKCEConfig struct {
	ClientID              string   // The client identifier to use.
	AuthorizationEndpoint string   // The URL of the authorization endpoint.
	TokenEndpoint         string   // The URL of the token endpoint.
	RedirectEndpoint      string   // The URL of the redirect endpoint.
	Scopes                []string // Recommended scope(s) to request.
}

// ClientCredentialsConfig contains OAuth 2.0 configuration for the Client Credentials Flow.
type ClientCredentialsConfig struct {
	TokenEndpoint string   // The URL of the token endpoint.
	Scopes        []string // Recommended scope(s) to request.
}

// OAuth2Config returns OAuth 2.0 configuration.
func (r Resolver) OAuth2Config() OAuth2Configuration {
	return r.c
}
