package muxer

import (
	"errors"
	"net/http"
	"regexp"
)

/*
Route defines a mapping between an HTTP request path and an HTTP request handler.
It contains the regular expression that matches the request path, the HTTP method,
the handler to be executed for that request, and the parameter names extracted from the path.
*/
type Route struct {
	path     *regexp.Regexp
	method   string
	handler  http.Handler
	params   []string
	template string
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

// PathTemplate retrieves the path template of the current route
func (r *Route) PathTemplate() (string, error) {
	if r == nil {
		return "", errors.New("route is nil, no template")
	}

	if r.template == "" {
		return r.template, errors.New("template is empty")
	}

	return r.template, nil
}
