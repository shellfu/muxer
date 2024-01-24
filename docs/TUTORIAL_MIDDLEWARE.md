# Using the `muxer` Package: Middleware

In the previous tutorial, we explored the basics of the `muxer` package,
which provides an HTTP request multiplexer for building web applications
in Go. We learned how to define routes and handle incoming requests. In
this tutorial, we will dive into the concept of middleware and see how it
can enhance the functionality of our web application.

## What is Middleware?

Middleware is a way to extend the functionality of an HTTP server by
adding additional processing steps to the request-response cycle. It sits
between the server and the application's handler functions, allowing you
to perform tasks such as authentication, logging, error handling, and
more.

Middleware functions receive an `http.Handler` as input and return a new
`http.Handler` that wraps the original handler. They can inspect and
modify the incoming request, perform pre-processing or post-processing
tasks, and decide whether to pass the request to the next middleware or
terminate the request-response cycle.

The `muxer` package provides a convenient way to use middleware with its
`Router` struct. Let's explore some built-in middleware provided by the
package and learn how to use them.

## CORS Middleware

CORS (Cross-Origin Resource Sharing) is a mechanism that allows web
browsers to make cross-origin requests. It involves adding special headers
to HTTP responses to indicate which origins are allowed to access the
resources on a web server.

The `muxer` package includes a CORS middleware that simplifies adding CORS
headers to your responses. To use the CORS middleware, follow these steps:

1. Import the `muxer` and `muxer/middleware` packages:
    ```go
    import (
        "net/http"

        "github.com/shellfu/muxer"
        "github.com/shellfu/muxer/middleware"
    )
    ```

2. Create a new `Router` instance:
    ```go
    router := muxer.NewRouter()
    ```

3. Create a new CORS middleware with default options:
    ```go
    cors := middleware.CORS()
    ```

4. Register your routes with the CORS middleware:
    ```go
    router.Use(cors)
    router.HandleFunc(http.MethodGet, "/api", myHandler)
    ```

Here, `myHandler` is the HTTP handler function that will be executed when the route is matched.

5. Start the server:
    ```go
    http.ListenAndServe(":8080", router)
    ```

With these steps, the CORS middleware will add the necessary CORS headers
to the responses of the `/api` route. By default, it allows any origin
(`*`), methods specified in the request, and headers specified in the
request.

You can customize the CORS middleware by passing `CORSOption` values to
its constructor. These options allow you to configure the allowed origins,
methods, headers, and other CORS settings according to your application's
requirements.

## Gzip Middleware

The Gzip middleware compresses the response body using gzip encoding if
the client supports it. This can significantly reduce the size of the
response and improve network performance. To use the Gzip middleware,
follow these steps:

1. Import the required packages:
    ```go
    import (
        "net/http"

        "github.com/shellfu/muxer"
        "github.com/shellfu/muxer/middleware"
    )
    ```

2. Create a new `Router` instance:
    ```go
    router := muxer.NewRouter()
    ```

3. Add the Gzip middleware to the router:
    ```go
    router.Use(middleware.Gzip())
    ```

The Gzip middleware will now compress the response body for all routes registered with the router.

4. Define your routes and handlers as usual:
    ```go
    router.HandleFunc("/", myHandler)
    ```

5. Start the server:
    ```go
    http.ListenAndServe(":8080", router)
    ```

Now, the Gzip middleware will automatically compress the response body
using gzip encoding if the client supports it. Otherwise, it will pass the
response to the next middleware or handler without modification.

## RecoveryHandler Middleware

The RecoveryHandler middleware is a useful tool for handling unexpected
panics in your application. It recovers from a panic, logs the panic
details, and sends an appropriate error response to the client. To use the
RecoveryHandler middleware, follow these steps:

1. Import the necessary packages:
    ```go
    import (
        "net/http"

        "github.com/shellfu/muxer"
        "github.com/shellfu/muxer/middleware"
    )
    ```

2. Create a new `Router` instance:
    ```go
    router := muxer.NewRouter()
    ```

3. Add the RecoveryHandler middleware to the router:
    ```go
    router.Use(middleware.RecoveryHandler(logger, true))
    ```

Here, `logger` is a custom logger that you can provide to log the panic
details. If `printStack` is set to `true`, the middleware will also log a
stack trace. If no logger is provided, the default Go logger will be used.

4. Define your routes and handlers as usual:
    ```go
    router.HandleFunc("/", myHandler)
    ```

5. Start the server:
    ```go
    http.ListenAndServe(":8080", router)
    ```

Now, if your application encounters a panic, the RecoveryHandler
middleware will catch it, log the details using the provided logger, and
return an appropriate HTTP 500 Internal Server Error response to the
client. This helps prevent your application from crashing due to
unexpected errors.

## Custom Middleware

The `muxer` package allows you to create your own custom middleware to
extend the functionality of your web application. Custom middleware
functions follow the same pattern as the built-in middleware functions:
they receive an `http.Handler` as input and return a new `http.Handler`
that wraps the original handler.

To create a custom middleware, follow these steps:

1. Define a middleware function with the following signature:
    ```go
    func myMiddleware(next http.Handler) http.Handler {
        // Perform pre-processing tasks here

        // Return a new handler that wraps the original handler
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Perform pre-processing tasks before calling the next handler

            // Call the next handler
            next.ServeHTTP(w, r)

            // Perform post-processing tasks after the next handler has been called
        })
    }
    ```

2. In the middleware function, you can perform any pre-processing tasks
   before calling the next handler, such as authentication, logging,
   modifying the request, or adding custom headers.

3. Create a new `Router` instance:
    ```go
    router := muxer.NewRouter()
    ```

4. Register your middleware function with the router using the `Use` method:
    ```go
    router.Use(myMiddleware)
    ```

5. Define your routes and handlers as usual:
    ```go
    router.HandleFunc("/", myHandler)
    ```

6. Start the server:
    ```go
    http.ListenAndServe(":8080", router)
    ```

With these steps, your custom middleware will be executed for each
incoming request, allowing you to add your own functionality to the
request-response cycle.

Custom middleware provides a flexible way to add functionality to your
application, allowing you to tailor it to your specific needs. You can
create multiple middleware functions and chain them together using the
`Use` method to execute them in a specific order.

Remember that the order in which you register your middleware matters.
Middleware functions registered earlier will be executed first, followed
by the ones registered later.

Feel free to experiment and create custom middleware that fits your
application requirements. It's a powerful tool for adding reusable and
modular functionality to your web application.

## Conclusion

In this tutorial, we explored the concept of middleware and learned how to
use middleware with the `muxer` package to enhance the functionality of
our web application. We covered the CORS middleware, which simplifies
adding CORS headers to responses, the Gzip middleware, which compresses
response bodies using gzip encoding, and the RecoveryHandler middleware,
which handles unexpected panics gracefully.

Using middleware can significantly improve the modularity,
maintainability, and extensibility of your Go web applications. It allows
you to separate concerns, add reusable functionality, and handle common
tasks without cluttering your route handlers.

By leveraging the middleware provided by the `muxer` package, you can
easily incorporate common functionality into your web application, saving
time and effort. Feel free to explore the built-in middleware further or
create your own custom middleware to suit your specific needs.

Happy coding!
