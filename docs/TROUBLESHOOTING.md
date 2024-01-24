# Troubleshooting for Muxer

This troubleshooting outline  provides solutions to common issues and errors 
that you may encounter when using the Muxer package. If you're experiencing
difficulties or unexpected behavior, follow the steps outlined below to
identify and resolve the problem.

## Table of Contents

1. [Missing Route Registration](#missing-route-registration)
2. [Route Not Found](#route-not-found)
3. [Incorrect Request Method](#incorrect-request-method)
4. [Parameter Extraction Failure](#parameter-extraction-failure)
5. [Middleware Execution Order](#middleware-execution-order)
6. [CORS Headers Not Set](#cors-headers-not-set)
7. [Gzip Compression Not Applied](#gzip-compression-not-applied)
8. [Recovery from Panic Not Working](#recovery-from-panic-not-working)

## 1. Missing Route Registration

### Problem
You're unable to access a specific endpoint or receive a "404 Not Found" error for a registered route.

### Solution
Make sure you have registered the route using the `Handle` or `HandleFunc` methods of the `Router` instance. Check the following:

- Verify that the route is registered with the correct HTTP method (`GET`, `POST`, etc.).
- Ensure that the path pattern used in the route registration matches the requested URL path.
- Double-check that the router instance is correctly created and used as the HTTP handler in your application.

## 2. Route Not Found

### Problem
All requests to your application result in a "404 Not Found" error, even though routes are registered correctly.

### Solution
Check the following:

- Ensure that the `ServeHTTP` method of your `Router` instance is correctly
  implemented and used as the HTTP handler for incoming requests. It should be
  called by the `http.ListenAndServe` function.
- Verify that the registered routes have the correct HTTP method and path patterns.
- If your application is using subrouters, ensure that they are correctly created and used.

## 3. Incorrect Request Method

### Problem
Requests to your application return a "405 Method Not Allowed" error, even though the request method is correct.

### Solution
Check the following:

- Make sure that the registered route has the correct HTTP method specified.
- Verify that the incoming request has the same HTTP method as the registered
  route. The comparison is case-sensitive, so ensure they match exactly (e.g.,`GET` vs. `get`).

## 4. Parameter Extraction Failure

### Problem
The parameters extracted from the URL path using the `Params` method of the `Router` instance are not populated correctly.

### Solution
Review the following:

- Ensure that the route registration pattern contains the named parameters using the `:paramName` syntax.
- Verify that the parameter names in the route registration pattern match the names used in the handler function or middleware.
- Check that the `Params` method is called with the correct `*http.Request` instance.

## 5. Middleware Execution Order

### Problem
The execution order of middleware functions is not as expected, leading to incorrect behavior or missing functionality.

### Solution
Check the following:

- Make sure the middleware functions are registered using the `Use` method of
  the `Router` instance in the correct order. Middleware functions are executed
  in the order they are registered.
- Verify that the middleware functions are applied to the correct router
  instance. If you're using subrouters, ensure that the middleware is registered
  on the appropriate router.

## 6. CORS Headers Not Set

### Problem
Cross-Origin Resource Sharing (CORS) headers are not set in the HTTP response, causing cross-origin requests to fail.

### Solution
Review the following:

- Ensure that the `CORS` middleware is correctly registered using the `Use` method of the `Router` instance.
- If you need to customize the CORS options, pass the appropriate `CORSOption` values when creating the `CORS` middleware.
- Check that the CORS middleware is registered before the routes that require CORS headers.

## 7. Gzip Compression Not Applied

### Problem
The response body is not compressed using Gzip encoding, even though the `Gzip` middleware is registered.

### Solution
Check the following:

- Make sure the `Gzip` middleware is correctly registered using the `Use` method of the `Router` instance.
- Ensure that the `Gzip` middleware is registered before the routes that require compression.
- Verify that the client supports Gzip encoding by sending an `Accept-Encoding` header with a value that includes `gzip`.

## 8. Recovery from Panic Not Working

### Problem
Your application panics, but the `RecoveryHandler` middleware doesn't recover from the panic or log any error messages.

### Solution
Check the following:

- Ensure that the `RecoveryHandler` middleware is correctly registered using the `Use` method of the `Router` instance.
- Make sure that the `RecoveryHandler` middleware is registered before the routes that may cause panics.
- If you're providing a custom logger, ensure that it implements the necessary logging functionality.

If the above steps don't resolve the issue, consider examining the
application's error logs and adding additional logging statements to identify
the cause of the panic.
