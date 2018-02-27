/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package graceful

import (
	"crypto/tls"
	"net/http"
	"time"
)

var (
	// DefaultCurvePreferences defines the recommended elliptic curves for modern TLS
	DefaultCurvePreferences = []tls.CurveID{
		tls.CurveP256,
		tls.X25519, // Go 1.8 only
	}

	// DefaultCipherSuites defines the recommended cipher suites for modern TLS
	DefaultCipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

		// Best disabled, as they don't provide Forward Secrecy,
		// but might be necessary for some clients
		// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}

	// DefaultMinVersion defines the recommended minimum version to use for the TLS protocol (1.2)
	DefaultMinVersion uint16 = tls.VersionTLS12

	// DefaultReadTimeout sets the maximum time a client has to fully stream a request (5s)
	DefaultReadTimeout = 5 * time.Second
	// DefaultWriteTimeout sets the maximum amount of time a handler has to fully process a request (10s)
	DefaultWriteTimeout = 10 * time.Second
	// DefaultIdleTimeout sets the maximum amount of time a Keep-Alive connection can remain idle before
	// being recycled (120s)
	DefaultIdleTimeout = 120 * time.Second
)

// WithDefaults patches a http.Server based on a best practice configuration
// from Cloudflare: https://blog.cloudflare.com/exposing-go-on-the-internet/
//
// You can override the defaults by mutating the Default* variables exposed
// by this package
func WithDefaults(srv *http.Server) *http.Server {
	if srv.TLSConfig == nil {
		srv.TLSConfig = &tls.Config{}
	}

	srv.TLSConfig.PreferServerCipherSuites = true
	srv.TLSConfig.MinVersion = DefaultMinVersion
	srv.TLSConfig.CurvePreferences = DefaultCurvePreferences
	srv.TLSConfig.CipherSuites = DefaultCipherSuites

	if srv.ReadTimeout == 0 {
		srv.ReadTimeout = DefaultReadTimeout
	}

	if srv.WriteTimeout == 0 {
		srv.WriteTimeout = DefaultWriteTimeout
	}

	if srv.IdleTimeout == 0 {
		srv.IdleTimeout = DefaultIdleTimeout
	}

	return srv
}
