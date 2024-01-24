# muxer - A lightweight HTTP router
[![Go Report Card](https://goreportcard.com/badge/github.com/shellfu/muxer)](https://goreportcard.com/report/github.com/shellfu/muxer)
[![GoDoc](https://godoc.org/github.com/shellfu/muxer?status.svg)](https://godoc.org/github.com/shellfu/muxer)
[![codecov](https://codecov.io/gh/shellfu/muxer/graph/badge.svg?token=2FFO3J3051)](https://codecov.io/gh/shellfu/muxer)

<p align="center"><img src="docs/img/logo.png" width="300"></p>

muxer is a lightweight HTTP router written in Go that is designed to be simple,
flexible, and familiar to users of the net/http package. It provides a simple
API for registering HTTP handlers and middleware functions, and matching incoming
requests to the appropriate handler based on their path and HTTP method.

muxer uses regular expressions to match incoming request paths, and allows for
named parameters to be extracted from the path and passed to the handler function.
It also supports middleware functions that can be executed before the main handler,
allowing for common functionality like authentication and logging to be shared
across multiple routes.

## Table of Contents
1. [Installation](#installation)
2. [Usage](#usage)
3. [Routing](#routing)
4. [Serving Static Files](#serving-static-files)
5. [Additional Documentation](#additional-documentation)

## Installation

To install muxer, use the following command:

```go
go get github.com/shellfu/muxer
```

Please note, if you encounter an issue where terminal prompts are disabled when using go.mod, you have two options to overcome this:
```
fatal: could not read Username for 'https://github.com': terminal prompts disabled
Confirm the import path was entered correctly.
If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.
```

### Option 1: Setting GOPRIVATE

You can configure the `GOPRIVATE` environment variable to bypass the prompts. `GOPRIVATE` is a comma-separated list of glob patterns (in the syntax of Go's path.Match) of module path prefixes.

Here is an example of how to set the `GOPRIVATE` environment variable:

```sh
export GOPRIVATE=github.com/shellfu/muxer
```

This setting will instruct Go to never download these modules from public proxies or check them into public checksum databases, and the modules will be fetched directly from the source.

### Option 2: Using a ~/.netrc file

Alternatively, you can create a ~/.netrc file to automatically authenticate with the git server when go get is run. This method requires you to have access credentials for the repository.

Here is an example of a ~/.netrc file:
```
machine github.com
login your-username
password your-password-or-personal-access-token
```
Remember to restrict the permissions of the ~/.netrc file to protect your credentials:


```sh
chmod 600 ~/.netrc
```

Make sure to replace your-username and your-password-or-personal-access-token with your actual GitHub username and either your password or a personal access token.

After setting up either `GOPRIVATE` or `~/.netrc`, you should be able to run go get github.com/shellfu/muxer without any issues.

## Usage

Here's an example of how to use muxer:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/shellfu/muxer"
)

func main() {
	r := muxer.NewRouter()

	r.HandleFunc(http.MethodGet, "/hello/:name", func(w http.ResponseWriter, r *http.Request) {
		params := router.Params(r)
		name := params["name"]
		fmt.Fprintf(w, "Hello, %s!", name)
	})

	http.ListenAndServe(":8080", r)
}
```

This example creates a new router using the `NewRouter` function, registers a
route using the `HandleFunc` method, and starts an HTTP server using the
`ListenAndServe` function.

Routes are defined using the `HandleFunc` method, which takes an HTTP method,
a path pattern, and an HTTP handler function. The path pattern may include named
parameters, indicated by a leading colon (e.g. "/users/:id"). The named parameters
are extracted from the path and added to the request context, where they can be
retrieved using the `Params` method of the Router.

muxer also supports middleware, which can be registered using the `Use` method of
the Router. Middleware functions take an `http.Handler` and return a new `http.Handler`
that performs some additional functionality before passing the request on to the
next handler in the chain.

For example, the following code adds a logger middleware that logs the HTTP request
method, URL, and duration:

```go
func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

r.Use(logger)
```

This middleware function takes an `http.Handler` and returns a new `http.Handler`
that logs the HTTP request method, URL, and duration before passing the request
on to the next handler in the chain.

muxer also provides several built-in middleware functions, such as CORS and Gzip
compression, which can be registered using the `Use` method and the corresponding
functions from the `middleware` package.

## Routing
The Router instance allows registering routes that match HTTP requests based on
the request's URL path and HTTP method. You can register a route using the `HandleRoute`
method, which takes an HTTP method, a path pattern, and an HTTP handler function.

To register a route with the Router instance, you can use the `HandleRoute` method,
which takes three arguments:

- An HTTP method (GET, POST, PUT, DELETE, etc.)
- A path pattern (e.g. /users/:id)
- An HTTP handler function that takes a http.ResponseWriter and a `*http.Request`

The path pattern can include named parameters that will be extracted from the URL
path and passed to the handler function as a map. For example, in the path pattern
/users/:id, :id is a named parameter that will match any string and be extracted
as a key-value pair in the map passed to the handler function.

Here's an example that shows how to register a route with the `Router` instance:

```go
r := muxer.NewRouter()

r.HandleRoute("GET", "/users/:id", func(w http.ResponseWriter, req *http.Request) {
    params := r.Params(req)
    id := params["id"]
    fmt.Fprintf(w, "User ID: %s", id)
})
```

In this example, we create a new Router instance and register a route that matches
GET requests to the path /users/:id. The route's handler function extracts the
value of the id parameter from the map returned by the `Params` method of the
Router instance and writes it to the response writer.

The HandleRoute method returns an error if there is a problem compiling the regular
expression that matches the path pattern. It's a good practice to handle this error
and log it appropriately.

In addition to the HandleRoute method, the muxer package also provides the `HandlerFunc`
method, which is a convenience method for registering a new route with an HTTP handler
function. It takes the same arguments as HandleRoute, but the HTTP handler function is
passed as an http.HandlerFunc instead of a generic http.Handler. This method is provided
to make the muxer API more familiar to users of the net/http package.

## Serving Static Files
While muxer doesn't provide built-in support for serving static files to avoid duplicating the functionality of http.FileServer, you can easily use http.FileServer in combination with muxer to handle static file serving. Here's how you can do it:

```go
router := muxer.NewRouter()

fs := http.FileServer(http.Dir("/path/to/your/static/files"))
router.Handle(http.MethodGet, "/static/:filepath", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	params := muxer.Params(req)
	req.URL.Path = params["filepath"]
	fs.ServeHTTP(w, req)
}))
```

In this example, an http.FileServer is created that serves files from a specified directory. Then a new route is registered in the muxer router that matches any GET requests where the path starts with "/static/". When such a request is matched, the remaining part of the path is used to identify the static file to be served.

This pattern can be adjusted to suit different needs, including serving files from multiple directories, serving files with certain content types, and more. The main idea here is to leverage the existing functionality of http.FileServer and integrate it with the routing capabilities of muxer to provide a seamless static file serving solution.

Remember to replace /path/to/your/static/files with the actual path to the directory that contains your static files.

Please note that it's always a good practice to handle any errors that may occur during the setup and usage of http.FileServer. This includes checking if the specified directory exists and is readable, and handling any errors returned by fs.ServeHTTP

---

## Additional Documentation

For more information about muxer, please refer to the following documents:

- [Basic Tutorial](docs/TUTORIAL_BASIC.md)
- [Middleware Tutorial](docs/TUTORIAL_MIDDLEWARE.md)
- [Prometheus Metrics Tutorial](docs/TUTORIAL_PROMETHEUS_METRICS.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
