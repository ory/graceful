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

	t.Run("case=timeout", func(t *testing.T) {
		server := WithDefaults(&http.Server{
			Addr:    "localhost:54932",
			Handler: &testServer{timeout: time.Second * 10},
		})

		// Start the server after 1s
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

		require.Error(t, <-done)
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
