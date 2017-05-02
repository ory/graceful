package graceful

import (
	"context"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"time"
	"net/http"
)

// HTTPServer is a wrapper for http.Server and makes running a graceful http server easy.
type HTTPServer struct {
	// ShutdownTimeout defines how long the server will wait before the shutdown is forced.
	ShutdownTimeout time.Duration

	stopChan        chan os.Signal
	*http.Server
}

// Graceful sets up graceful shutdown for the HTTPServer, or returns an error if something goes wrong.
//
//   if err := server.Graceful(server, func () {
//   	   if err := server.ListenAndServe(); err != nil {
//		   // ...
//	   }
//   }); err != nil {
//   	   // ..
//   }
func (srv *HTTPServer) Graceful(runner func()) error {
	signal.Notify(srv.stopChan, os.Interrupt)

	go runner()

	<-srv.stopChan // wait for SIGINT
	timer, cancel := context.WithTimeout(context.Background(), srv.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(timer); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
