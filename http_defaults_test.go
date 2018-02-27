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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDefaults(t *testing.T) {
	for k, tc := range []*http.Server{
		{},
		{
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
