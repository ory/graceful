package graceful

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"
)

var cfCurves = []tls.CurveID{
	tls.CurveP256,
	tls.X25519, // Go 1.8 only
}

var cfCipherSuites = []uint16{
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

var cfTLSConfig = &tls.Config{
	PreferServerCipherSuites: true,
	CurvePreferences:         cfCurves,
	MinVersion:               tls.VersionTLS12,
	CipherSuites:             cfCipherSuites,
}

// PatchHTTPServerWithCloudflareConfig patches a http.Server based on a best practice configuration
// from Cloudflare: https://blog.cloudflare.com/exposing-go-on-the-internet/
func PatchHTTPServerWithCloudflareConfig(srv *http.Server) *HTTPServer {
	if srv.TLSConfig == nil {
		srv.TLSConfig = cfTLSConfig
	}

	srv.TLSConfig.PreferServerCipherSuites = true
	srv.TLSConfig.CurvePreferences = cfCurves
	srv.TLSConfig.MinVersion = tls.VersionTLS12
	srv.TLSConfig.CipherSuites = cfCipherSuites

	if srv.ReadTimeout == 0 {
		srv.ReadTimeout = time.Second * 5
	}

	if srv.WriteTimeout == 0 {
		srv.WriteTimeout = time.Second * 10
	}

	if srv.IdleTimeout == 0 {
		srv.IdleTimeout = time.Second * 120
	}

	return &HTTPServer{
		Server:          srv,
		ShutdownTimeout: 5 * time.Second,
		stopChan:        make(chan os.Signal),
	}
}
