package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

/*
Gzip is a middleware function that returns a new HTTP handler function
that compresses the response body using gzip encoding if the client accepts it.
It modifies the response headers to include the 'Content-Encoding' and 'Vary'
headers, and wraps the response writer with a gzip writer to compress the body.

If the client doesn't support gzip encoding, it just calls the next handler
in the chain without modifying the response.

Example usage:

r := muxer.NewRouter()
r.Use(muxer.Gzip())
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
 fmt.Fprintln(w, "Hello, world!")
})

http.ListenAndServe(":8080", r)
*/
func Gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		handler.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

// A gzipResponseWriter wraps an http.ResponseWriter and a gzip.Writer
// to compress the response.
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
