package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/pkg/errors"
)

func GracefullyRunHTTPServer(srv *http.Server, runner func()) error {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	go runner()

	<-stopChan // wait for SIGINT
	timer, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := srv.Shutdown(timer); err != nil {
		errors.WithStack(err)
	}

	return nil
}
