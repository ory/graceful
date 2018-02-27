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

package graceful_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/ory/graceful"
)

func ExampleGraceful() {
	server := graceful.WithDefaults(&http.Server{
		Addr: "localhost:8080",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			fmt.Println("handler: Received the request")
			time.Sleep(3 * time.Second)

			fmt.Println("handler: Fulfilling the request after 3 seconds")
			fmt.Fprint(rw, "Hello World!")
		}),
	})

	// Kill the server after 5 seconds
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("human: Killing the server after 2 seconds")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	// Start the server
	done := make(chan struct{})
	go func() {
		fmt.Println("graceful: Starting the server")
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			fmt.Println("graceful: Failed to gracefully shutdown")
			os.Exit(-1)
		}
		fmt.Println("graceful: Server was shutdown gracefully")

		done <- struct{}{}
	}()

	time.Sleep(1 * time.Second) // Give the server time to start up

	fmt.Println("main: Sending request")
	res, _ := http.Get("http://localhost:8080/")
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("main: Received response ->", string(body))

	<-done

	// Output:
	// graceful: Starting the server
	// main: Sending request
	// handler: Received the request
	// human: Killing the server after 2 seconds
	// handler: Fulfilling the request after 3 seconds
	// main: Received response -> Hello World!
	// graceful: Server was shutdown gracefully
}
