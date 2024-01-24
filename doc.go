/*
Package muxer provides an HTTP request multiplexer that allows you to register multiple HTTP handlers
for different paths and methods on a single server. It is designed to be lightweight, flexible, and
easy to use, making it a popular choice for building web applications and APIs in Go.

The Router type implements the http.Handler interface and can be used as the server's primary
request handler. It provides several methods for registering routes, middleware, and error handlers,
as well as extracting named parameters from the URL path and handling unsupported HTTP methods.

Route registration is done using the Handle, HandleFunc, and HandleRoute methods, which allow you
to specify the HTTP method, URL path, and handler function for each route. Path parameters are supported
using named placeholders in the path, such as "/users/:id". You can extract these parameters from the
request using the Params method of the Router type, which returns a map of parameter names to their
corresponding values.

Middleware functions can also be registered using the Use method, which allows you to chain multiple
middleware functions together in a specific order. Middleware functions are executed before the main
handler function, and can be used to perform tasks such as authentication, logging, or request/response
processing.

The Router type also supports error handling using the NotFoundHandler and PanicHandler fields. The
NotFoundHandler is executed when a request is made for a path that does not match any registered route,
and returns a 404 Not Found HTTP status code. The PanicHandler is executed when a panic occurs during
route processing, and can be used to handle and recover from unexpected errors.

	Example usage:

	  package main

	  import (
	    "fmt"
	    "log"
	    "net/http"

	    "github.com/shellfu/muxer"
	  )

	  func main() {
	    // create a new Router instance
	    router := muxer.NewRouter()

	    // register a route using the Handle method
	    router.Handle("GET", "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	      fmt.Fprint(w, "Hello, world!")
	    }))

	    // register a route using the HandleFunc method
	    router.HandleFunc("GET", "/users/:id", func(w http.ResponseWriter, r *http.Request) {
	      // extract the "id" parameter from the URL path using Params()
	      params := router.Params(r)
	      id := params["id"]

	      // handle the request
	      // ...
	    })

	    // register a route using the HandleRoute method
	    router.HandleRoute("POST", "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	      // extract the "id" parameter from the URL path using Params()
	      params := router.Params(r)
	      id := params["id"]

	      // handle the request
	      // ...
	    }))

	    // register middleware functions using the Use method
	    router.Use(func(next http.Handler) http.Handler {
	      return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	        // do some preprocessing before passing the request to the next handler
	        // ...
	        next.ServeHTTP(w, r)
	      })
	    })

	    // start the HTTP server and listen for incoming requests
	    log.Fatal(http.ListenAndServe(":8080", router))
	  }

	  For more information and examples, see the documentation for individual methods and types.
*/
package muxer
