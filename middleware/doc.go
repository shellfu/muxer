/*
Package middleware provides HTTP middleware handlers.

CORS middleware adds Cross-Origin Resource Sharing (CORS) headers to the HTTP response. By default, it sets the Access-Control-Allow-Origin header to "*", which allows any origin to access the resource. It also sets the Access-Control-Allow-Methods header to the HTTP methods defined in the Config, and the Access-Control-Allow-Headers header to the HTTP headers defined in the Config.

The middleware can be customized by passing in one or more CORSOption values to the constructor. These options can be used to configure the allowed origins, methods, headers, and other CORS settings.

Usage:

		// Create a new Router
		router := muxer.NewRouter()

		// Create a new CORS middleware with default options
		cors := muxer.CORS()

		// Register a new route with the CORS middleware
		router.HandleFunc("/api", myHandler).Methods("GET").Middleware(cors)

		// Start the server
		log.Fatal(http.ListenAndServe(":8080", router))

	 -------------------------------------------------------------------------

Gzip middleware returns a new HTTP handler function that compresses the response body using gzip encoding if the client accepts it. It modifies the response headers to include the 'Content-Encoding' and 'Vary' headers, and wraps the response writer with a gzip writer to compress the body.

If the client doesn't support gzip encoding, it just calls the next handler in the chain without modifying the response.

Example usage:

	 r := muxer.NewRouter()
	 r.Use(muxer.Gzip())
	 r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		 fmt.Fprintln(w, "Hello, world!")
	 })

	 http.ListenAndServe(":8080", r)

	 -------------------------------------------------------------------------

RecoveryHandler middleware is an HTTP middleware that recovers from a panic, logs the panic, writes http.StatusInternalServerError, and continues to the next handler.

Usage:

	r := muxer.NewRouter()
	r.Use(middleware.RecoveryHandler(
		myCustomLogger{},
		true,
	))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("Unexpected error!")
	})
	http.ListenAndServe(":1123", r)

The RecoveryHandler logs errors and, if printStack is true, also logs a stack trace. If printStack is false, no stack trace is logged. If no logger is provided, it uses the default Go logger.
*/
package middleware
