package muxer

import (
	"net/http"
)

/*
A RouterOption is a function that sets a configuration option on a Router.
It takes a reference to a Router and modifies its properties.
*/
type RouterOption func(r *Router)

/*
WithNotFoundHandler option takes a http.Handler that will be set as the
NotFoundHandler of the Router. This handler will be executed when the
Router receives a request for an unknown path.
*/
func WithNotFoundHandler(handler http.Handler) RouterOption {
	return func(r *Router) {
		r.NotFoundHandler = handler.(http.HandlerFunc)
	}
}

/*
WithMaxRequestBodySize option sets the maximum size of the request body that
the Router can handle. This option can be used to prevent denial-of-service
attacks by limiting the size of the request body.
*/
func WithMaxRequestBodySize(size int64) RouterOption {
	return func(r *Router) {
		r.MaxRequestBodySize = size
	}
}
