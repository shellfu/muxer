package muxer

import (
	"net/http"
	"regexp"
)

/*
Route defines a mapping between an HTTP request path and an HTTP request handler.
It contains the regular expression that matches the request path, the HTTP method,
the handler to be executed for that request, and the parameter names extracted from the path.
*/
type Route struct {
	path    *regexp.Regexp
	method  string
	handler http.Handler
	params  []string
}

func (r *Route) match(path string) map[string]string {
	match := r.path.FindStringSubmatch(path)
	if match == nil {
		return nil
	}

	params := make(map[string]string)
	for i, name := range r.params {
		params[name] = match[i+1]
	}

	return params
}
