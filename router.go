package muxer

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type contextKey string

const (
	// ParamsKey is the key used to store the extracted parameters in the request context.
	ParamsKey contextKey = "params"
	// RouteContextKey is the key used to store the matched route in the request context
	RouteContextKey contextKey = "matched_route"
)

/*
Router is an HTTP request multiplexer. It contains the registered routes and middleware functions.
It implements the http.Handler interface to be used with the http.ListenAndServe function.
*/
type Router struct {
	http.Handler

	routes     []Route
	middleware []func(http.Handler) http.Handler
	subrouters map[string]*Router

	NotFoundHandler    http.HandlerFunc
	MaxRequestBodySize int64
}

// NewRouter creates a new instance of a Router with optional configuration provided
// through the RouterOptions
func NewRouter(options ...RouterOption) *Router {
	r := &Router{
		NotFoundHandler: http.HandlerFunc(http.NotFound),
		subrouters:      make(map[string]*Router),
	}

	for _, option := range options {
		option(r)
	}

	return r
}

/*
Subrouter returns a new router that will handle requests that match the given attribute value.
The attribute value can be, for example, a host or path prefix. If a subrouter does not already exist
for the given attribute value, a new one will be created. The new router will inherit the parent router's
NotFoundHandler and other settings.
*/
func (r *Router) Subrouter(attrValue string) *Router {
	if _, ok := r.subrouters[attrValue]; !ok {
		// If subrouter doesn't exist for attribute value, create one
		subrouter := &Router{
			NotFoundHandler: r.NotFoundHandler,
			middleware:      append([]func(http.Handler) http.Handler{}, r.middleware...),
			subrouters:      make(map[string]*Router),
		}
		r.subrouters[attrValue] = subrouter
	}
	return r.subrouters[attrValue]
}

/*
Handle registers a new route with the given method, path and handler.

The method parameter specifies the HTTP method (e.g. GET, POST, PUT, DELETE, etc.) that
the route should match. If an unsupported method is passed, an error will be returned.

The path parameter specifies the URL path that the route should match. Path parameters
are denoted by a colon followed by the parameter name (e.g. "/users/:id").

The handler parameter is the HTTP handler function that will be executed when the route
is matched. The handler function should take an http.ResponseWriter and an *http.Request
as its parameters.
*/
func (r *Router) Handle(method string, path string, handler http.Handler) {
	r.HandlerFunc(method, path, func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(w, req)
	})
}

/*
HandlerFunc is a convenience method for registering a new route with an HTTP handler function.
It is similar to the net/http.HandleFunc method, and is provided to make the Router API more familiar
to users of the net/http package.

The method takes an HTTP method, a path pattern, and an HTTP handler function. The path pattern may
include named parameters, indicated by a leading colon (e.g. "/users/:id"). The named parameters are
extracted from the path and added to the request context, where they can be retrieved using the Params
method of the Router.

The handler function may be provided as an http.HandlerFunc, or as any other function that satisfies
the http.Handler interface (e.g. a method of a struct that implements ServeHTTP).
*/
func (r *Router) HandlerFunc(method, path string, handlerFunc http.HandlerFunc) {
	r.HandleRoute(method, path, handlerFunc)
}

/*
HandleRoute registers a new route with the given HTTP method, path, and handler function.
It adds the route to the router's list of routes and extracts the named parameters from the path
using regular expressions.

The method parameter specifies the HTTP method (e.g. GET, POST, PUT, DELETE, etc.) that
the route should match. If an unsupported method is passed, an error will be returned.

The path parameter specifies the URL path that the route should match. Path parameters
are denoted by a colon followed by the parameter name (e.g. "/users/:id").

The handler parameter is the HTTP handler function that will be executed when the route
is matched. The handler function should take an http.ResponseWriter and an *http.Request
as its parameters.

	Example usage:
	  router := muxer.NewRouter()
	  router.HandleRoute("GET", "/users/:id", func(w http.ResponseWriter, r *http.Request) {
	      // extract the "id" parameter from the URL path using Params()
	      params := router.Params(r)
	      id := params["id"]

	      // handle the request
	      // ...
	  })

	If there's an error compiling the regular expression that matches the path, it returns the error.
*/
func (r *Router) HandleRoute(method, path string, handler http.HandlerFunc) {
	// Parse path to extract parameter names
	paramNames := make([]string, 0)
	re := regexp.MustCompile(`:([\w-]+)`)
	pathRegex := re.ReplaceAllStringFunc(path, func(m string) string {
		paramName := m[1:]
		paramNames = append(paramNames, paramName)
		return `([-\w.]+)`
	})

	exactPath := regexp.MustCompile("^" + pathRegex + "$")

	r.routes = append(r.routes, Route{
		method:   method,
		path:     exactPath,
		handler:  handler,
		params:   paramNames,
		template: path, // Save the original template
	})
}

// HandlerFuncWithMethods is a convenience method for registering a new route with multiple HTTP methods.
// It is similar to the net/http.HandleFunc method, and is provided to make the Router API more familiar
// to users of the net/http package.
//
// The method takes a slice of HTTP methods, a path pattern, and an HTTP handler function. The path pattern may
// include named parameters, indicated by a leading colon (e.g. "/users/:id"). The named parameters are
// extracted from the path and added to the request context, where they can be retrieved using the Params
// method of the Router.
func (r *Router) HandlerFuncWithMethods(methods []string, path string, handlerFunc http.HandlerFunc) {
	for _, method := range methods {
		r.HandlerFunc(method, path, handlerFunc)
	}
}

/*
ServeHTTP dispatches the HTTP request to the registered handler that matches
the HTTP method and path of the request. It executes the middleware functions
in reverse order and sets the extracted parameters in the request context.
If there's no registered route that matches the request, it returns a
404 HTTP status code.
*/
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.MaxRequestBodySize > 0 && req.Body != nil {
		if req.ContentLength <= r.MaxRequestBodySize {
			req.Body = http.MaxBytesReader(w, req.Body, r.MaxRequestBodySize)
		} else {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
	}

	// Check subrouters first
	for prefix, subrouter := range r.subrouters {
		var matched bool
		switch {
		case prefix == req.URL.Host:
			matched = true
		case strings.HasPrefix(req.URL.Path, prefix):
			matched = true
			req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		}

		if matched {
			subrouter.ServeHTTP(w, req)
			return
		}
	}

	var methodMismatch bool
	for _, route := range r.routes {
		if route.method != req.Method {
			methodMismatch = true
			continue
		}
		params := route.match(req.URL.Path)
		if params == nil {
			continue
		}

		ctx := req.Context()
		ctx = context.WithValue(ctx, ParamsKey, params)
		ctx = context.WithValue(ctx, RouteContextKey, &route)

		handler := route.handler
		for i := len(r.middleware) - 1; i >= 0; i-- {
			handler = r.middleware[i](handler)
		}

		handler.ServeHTTP(w, req.WithContext(ctx))
		return
	}

	if methodMismatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.NotFoundHandler.ServeHTTP(w, req)
}

/*
Params returns the parameter names and values extracted from the request path.
It extracts the parameters from the request context, returns an empty map if
there are no parameters found.
*/
func (r *Router) Params(req *http.Request) map[string]string {
	return Params(req)
}

/*
Params returns the parameter names and values extracted from the request path.
It extracts the parameters from the request context, returns an empty map if
there are no parameters found.
*/
func Params(req *http.Request) map[string]string {
	params := req.Context().Value(ParamsKey)
	if p, ok := params.(map[string]string); ok {
		return p
	}
	return make(map[string]string)
}

/*
Use registers middleware functions that will be executed before the main handler.
It chains the middleware functions to create a new handler that executes them in
the given order before executing the main handler.
*/
func (r *Router) Use(middleware ...func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, middleware...)
}

// CurrentRoute returns the matched route for the current request, if any.
// This only works when called inside the handler of the matched route
// because the matched route is stored inside the request's context,
// which is wiped after the handler returns.
func CurrentRoute(r *http.Request) *Route {
	if rv := r.Context().Value(RouteContextKey); rv != nil {
		return rv.(*Route)
	}
	return nil
}
