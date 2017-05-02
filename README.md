# graceful

[![Build Status](https://travis-ci.org/ory/graceful.svg?branch=master)](https://travis-ci.org/ory/graceful)
[![Coverage Status](https://coveralls.io/repos/github/ory/graceful/badge.svg?branch=master)](https://coveralls.io/github/ory/graceful?branch=master)

Best practice http server configurations and helpers for Go 1.8's http graceful shutdown feature. Currently supports
best practice configurations by:

* [Cloudflare](https://blog.cloudflare.com/exposing-go-on-the-internet/)

## Usage

To install this library, do:

```sh
go get github.com/ory/graceful
```

### Running Cloudflare Config with Graceful Shutdown

```go
package main

import "github.com/ory/graceful"
import "net/http"

func main() {
    server := graceful.PatchHTTPServerWithCloudflareConfig(&http.Server{
        Addr: ":54932",
        // Handler: someHandler,
    })

    if err := server.Graceful(func () {
        if err := server.ListenAndServe(); err != nil {
            // ...
        }
    }); err != nil {
        // ...
    }
}
```
