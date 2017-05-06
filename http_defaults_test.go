package graceful

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDefaults(t *testing.T) {
	for k, tc := range []*http.Server{
		&http.Server{},
		&http.Server{
			TLSConfig: &tls.Config{
				PreferServerCipherSuites: false,
				CurvePreferences:         []tls.CurveID{},
				MinVersion:               tls.VersionTLS10,
				CipherSuites:             []uint16{},
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			server := WithDefaults(tc)
			assert.Equal(t, DefaultCurvePreferences, server.TLSConfig.CurvePreferences)
			assert.Equal(t, DefaultCipherSuites, server.TLSConfig.CipherSuites)
			assert.Equal(t, DefaultMinVersion, server.TLSConfig.MinVersion)
			assert.True(t, server.IdleTimeout > 0)
			assert.True(t, server.ReadTimeout > 0)
			assert.True(t, server.WriteTimeout > 0)
		})
	}
}
