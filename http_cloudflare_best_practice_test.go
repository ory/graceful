package graceful

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/tls"
	"fmt"
)

func TestPatchHTTPServerWithCloudflareConfig(t *testing.T) {
	for k, tc := range []*http.Server{
		&http.Server{},
		&http.Server{
			TLSConfig: &tls.Config{
				PreferServerCipherSuites: false,
				CurvePreferences:  []tls.CurveID{},
				MinVersion: tls.VersionTLS10,
				CipherSuites: []uint16{},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			server := PatchHTTPServerWithCloudflareConfig(tc)
			assert.Equal(t, cfCurves, server.TLSConfig.CurvePreferences)
			assert.Equal(t, cfCipherSuites, server.TLSConfig.CipherSuites)
			assert.Equal(t, uint16(tls.VersionTLS12), server.TLSConfig.MinVersion)
			assert.True(t, server.IdleTimeout > 0)
			assert.True(t, server.ReadTimeout > 0)
			assert.True(t, server.WriteTimeout > 0)
		})
	}
}