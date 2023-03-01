// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package graceful

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testServer struct {
	delay time.Duration
}

func (s *testServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	time.Sleep(s.delay)
	rw.Write([]byte("hi"))
}

func TestGraceful(t *testing.T) {
	t.Run("case=in-time", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54931",
			Handler: &testServer{delay: time.Second * 3},
		})

		go func() {
			require.NoError(t, Graceful(server.ListenAndServe, server.Shutdown))
		}()

		res, err := http.Get("http://localhost:54931/")

		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		all, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte("hi"), all)
	})

	t.Run("case=timeout during shutdown", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54932",
			Handler: &testServer{delay: time.Second * 10},
		})

		done := make(chan error)
		go func() {
			done <- Graceful(server.ListenAndServe, server.Shutdown)
		}()

		// Kill the server after 1s
		go func() {
			time.Sleep(1 * time.Second)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()

		_, err := http.Get("http://localhost:54932/")
		require.Error(t, err)

		require.ErrorIs(t, <-done, context.DeadlineExceeded)
	})

	t.Run("case=canceled", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54933",
			Handler: &testServer{delay: 0},
		})

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error)
		go func() {
			done <- GracefulContext(ctx, server.ListenAndServe, server.Shutdown)
		}()

		_, err := http.Get("http://localhost:54933/")
		require.NoError(t, err)
		cancel()
		require.NoError(t, <-done)
	})

	t.Run("case=canceled (timeout during shutdown)", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54934",
			Handler: &testServer{delay: time.Second * 10},
		})

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error)
		go func() {
			done <- GracefulContext(ctx, server.ListenAndServe, server.Shutdown)
		}()

		// Kill the server after 1s
		go func() {
			time.Sleep(1 * time.Second)
			cancel()
		}()

		_, err := http.Get("http://localhost:54934/")
		require.Error(t, err)

		require.ErrorIs(t, <-done, context.DeadlineExceeded)
	})

	t.Run("case=start-error", func(t *testing.T) {
		startErr := errors.New("Test error")

		start := func() error { return startErr }
		shutdown := func(c context.Context) error {
			return nil
		}

		err := Graceful(start, shutdown)
		require.Equal(t, startErr, err)
	})

	time.Sleep(time.Second) // clean up
}
