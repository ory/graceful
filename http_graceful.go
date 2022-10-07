// Copyright © 2022 Ory Corp

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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// StarFunc is the type of the function invoked by Graceful to start the server
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
	var (
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)

	// Setup the graceful shutdown handler (traps SIGINT and SIGTERM)
	go func() {
		signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		<-stopChan

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
