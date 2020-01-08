// Copyright (c) 2020, Sylabs, Inc. All rights reserved.

package server

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	code int
	size int
}

// WriteHeader records the HTTP status code as it is written.
func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.code = code
	lw.ResponseWriter.WriteHeader(code)
}

// Write accumulates the response size as the response is written.
func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lw.ResponseWriter.Write(b)
	lw.size += n
	return n, err
}

// remoteIP attempts to find the remote IP associated with a HTTP request.
func remoteIP(req *http.Request) string {
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ", ")[0]
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	} else if ip, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		return ip
	}
	return ""
}

// loggingHandler logs details about a HTTP request.
func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &loggingResponseWriter{w, http.StatusOK, 0}
		next.ServeHTTP(lw, r)

		entry := logrus.WithFields(logrus.Fields{
			"remote":  remoteIP(r),
			"host":    r.Host,
			"method":  r.Method,
			"path":    r.RequestURI,
			"referer": r.Referer(),
			"agent":   r.UserAgent(),
			"code":    lw.code,
			"size":    lw.size,
			"took":    time.Since(start),
		})
		entry.Info("completed handling request")
	})
}
