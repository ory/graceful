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
	"context"
	"io"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testServer struct {
	timeout time.Duration
}

func (s *testServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	time.Sleep(s.timeout)
	rw.Write([]byte("hi"))
}

func TestGraceful(t *testing.T) {
	t.Run("case=in-time", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54931",
			Handler: &testServer{timeout: time.Second * 3},
		})

		go func() {
			require.NoError(t, Graceful(server))
		}()

		res, err := http.Get("http://localhost:54931/")

		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		all, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, []byte("hi"), all)
	})

	t.Run("case=timeout", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54932",
			Handler: &testServer{timeout: time.Second * 10},
		})

		// Start the server
		done := make(chan error)
		go func() {
			done <- Graceful(server)
		}()

		// Kill the server after 1s
		go func() {
			time.Sleep(1 * time.Second)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()

		_, err := http.Get("http://localhost:54932/")
		require.Error(t, err)

		require.Error(t, <-done)
	})

	t.Run("case=shutdown", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr: "localhost:54933",
		})

		// Shutdown the server after 1s
		go func() {
			time.Sleep(1 * time.Second)
			err := server.Shutdown(context.Background())
			require.NoError(t, err)
		}()

		err := Graceful(server)
		require.NoError(t, err)
	})

	time.Sleep(time.Second) // clean up
}
