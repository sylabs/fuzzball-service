// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sylabs/compute-service/internal/pkg/model"
	"gopkg.in/square/go-jose.v2"
)

// discoverAuthMetadata attempts to discover metadata from an OAuth issuer using well-known URI
// discovery for OAuth 2.0 (RFC 8414) and OpenID Connect (OpenID.Discovery).
func discoverAuthMetadata(ctx context.Context, hc *http.Client, uri string) (model.AuthMetadata, error) {
	md, err := getAuthMetadata(ctx, hc, uri+"/.well-known/oauth-authorization-server")
	if err != nil {
		if md, err = getAuthMetadata(ctx, hc, uri+"/.well-known/openid-configuration"); err != nil {
			return model.AuthMetadata{}, err
		}
	}
	return md, nil
}

// getAuthMetadata gets metadata from uri as per the OAuth 2.0 Authorization Server Metadata
// specification (RFC 8414).
func getAuthMetadata(ctx context.Context, hc *http.Client, uri string) (md model.AuthMetadata, err error) {
	logrus.Info("getting auth metadata")
	defer func(t time.Time) {
		log := logrus.WithField("took", time.Since(t))
		if err != nil {
			log.WithError(err).Warn("failed to get auth metadata")
		} else {
			log.WithField("metadata", fmt.Sprintf("%#v", md)).Info("got auth metadata")
		}
	}(time.Now())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return model.AuthMetadata{}, err
	}
	res, err := hc.Do(req)
	if err != nil {
		return model.AuthMetadata{}, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&md); err != nil {
		return model.AuthMetadata{}, err
	}
	return md, nil
}

// getKeySet gets a JSON Web Key Set from uri as per the JSON Web Key specification (RFC 7515).
func getKeySet(ctx context.Context, hc *http.Client, uri string) (ks jose.JSONWebKeySet, err error) {
	logrus.Info("getting key set")
	defer func(t time.Time) {
		log := logrus.WithField("took", time.Since(t))
		if err != nil {
			log.WithError(err).Warn("failed to get key set")
		} else {
			log.WithField("keySet", fmt.Sprintf("%+v", ks)).Info("got key set")
		}
	}(time.Now())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	res, err := hc.Do(req)
	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&ks); err != nil {
		return jose.JSONWebKeySet{}, err
	}
	return ks, nil
}
