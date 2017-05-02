package graceful

import (
	"io/ioutil"
	"net/http"
	"os"
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

func TestGracefullyRunHTTPServer(t *testing.T) {
	t.Run("case=in-time", func(t *testing.T) {
		server := PatchHTTPServerWithCloudflareConfig(&http.Server{
			Addr:    ":54931",
			Handler: &testServer{timeout: time.Second * 3},
		})

		go func() {
			require.NoError(t, server.Graceful(func() {
				server.ListenAndServe()
			}))
		}()

		res, err := http.Get("http://127.0.0.1:54931/")

		server.stopChan <- os.Interrupt

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		all, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		assert.Equal(t, []byte("hi"), all)

	})

	t.Run("case=timeout", func(t *testing.T) {
		server := PatchHTTPServerWithCloudflareConfig(&http.Server{
			Addr:    ":54932",
			Handler: &testServer{timeout: time.Second * 10},
		})

		go func() {
			require.NoError(t, server.Graceful(func() {
				server.ListenAndServe()
			}))
		}()

		_, err := http.Get("http://127.0.0.1:54932/")
		server.stopChan <- os.Interrupt
		require.Error(t, err)
	})

	time.Sleep(time.Second) // clean up
}
