package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// RecoveryLogger is an interface used by the RecoveryHandler to log errors.
type RecoveryLogger interface {
	Println(v ...interface{})
}

// recoveryHandler is an HTTP middleware that recovers from a panic, logs the panic,
// writes http.StatusInternalServerError, and continues to the next handler.
type recoveryHandler struct {
	handler    http.Handler
	logger     RecoveryLogger
	printStack bool
}

/*
RecoveryHandler is a middleware that recovers from a panic, logs the panic,
writes http.StatusInternalServerError, and continues to the next handler.

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

The RecoveryHandler logs errors and, if printStack is true, also logs a
stack trace. If printStack is false, no stack trace is logged. If no logger is
provided, it uses the default Go logger.
*/
func RecoveryHandler(logger RecoveryLogger, printStack bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &recoveryHandler{handler: next, logger: logger, printStack: printStack}
	}
}

func (rh *recoveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rh.log(err)
		}
	}()

	rh.handler.ServeHTTP(w, r)
}

func (rh *recoveryHandler) log(v ...interface{}) {
	if rh.logger != nil {
		rh.logger.Println(v...)
	} else {
		if len(v) > 0 {
			log.Println(v[0])
		}
	}

	if rh.printStack {
		stack := string(debug.Stack())
		if rh.logger != nil {
			rh.logger.Println(stack)
		} else {
			log.Println(stack)
		}
	}
}
