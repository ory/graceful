// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// StartFunc is the type of the function invoked by Graceful to start the server
type StartFunc func() error

// ShutdownFunc is the type of the function invoked by Graceful to shutdown the server
type ShutdownFunc func(context.Context) error

// DefaultShutdownTimeout defines how long Graceful will wait before forcibly shutting down
var DefaultShutdownTimeout = 5 * time.Second

// Graceful sets up graceful handling of SIGINT and SIGTERM, typically for an HTTP server.
// When signal is trapped, the shutdown handler will be invoked with a context that expires
// after DefaultShutdownTimeout (5s).
//
//	server := graceful.WithDefaults(http.Server{})
//
//	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
//		log.Fatal("Failed to gracefully shut down")
//	}
func Graceful(start StartFunc, shutdown ShutdownFunc) error {
	return GracefulContext(context.Background(), start, shutdown)
}

// GracefulContext works like Graceful, but also shuts down the server when ctx
// is done, i.e., when the channel returned from ctx.Done() is closed. Any error
// from ctx.Err() is discarded.
func GracefulContext(ctx context.Context, start StartFunc, shutdown ShutdownFunc) error {
	var (
		stopChan = make(chan os.Signal, 1)
		errChan  = make(chan error, 1)
	)

	// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
	go func() {
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-stopChan:
		case <-ctx.Done():
		}

		timer, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		if err := shutdown(timer); err != nil {
			errChan <- errors.WithStack(err)
			return
		}

		errChan <- nil
	}()

	// Start the server
	if err := start(); err != http.ErrServerClosed {
		return err
	}

	return <-errChan
}
