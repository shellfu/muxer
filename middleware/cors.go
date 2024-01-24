package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

/*
The corsConfig struct contains the allowed origins, methods, and headers for the CORS.
The AllowCredentials field is used to allow or deny sending credentials such as cookies
or HTTP authentication. The MaxAge field is used to set the maximum age of the preflight
request cache.
*/
type corsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   map[string]string
	PreflightHeaders map[string]string
	MaxAge           int
}

// CORSOption is a function that modifies the CORSConfig.
type CORSOption func(*corsConfig)

// WithAllowedOrigins sets the list of allowed origins in the CORSConfig.
func WithAllowedOrigins(origins ...string) CORSOption {
	return func(cfg *corsConfig) {
		cfg.AllowedOrigins = origins
	}
}

// WithAllowedMethods sets the list of allowed methods in the CORSConfig.
func WithAllowedMethods(methods ...string) CORSOption {
	return func(cfg *corsConfig) {
		cfg.AllowedMethods = methods
	}
}

// WithAllowedHeaders sets the list of allowed headers in the CORSConfig.
// This function takes a slice of strings representing the allowed headers.
// It creates a new map with the header names as keys and empty string values.
// The map is then set as the AllowedHeaders field in the corsConfig struct.
func WithAllowedHeaders(headers ...string) CORSOption {
	headerMap := make(map[string]string, len(headers))
	for _, header := range headers {
		headerMap[header] = ""
	}
	return func(cfg *corsConfig) {
		cfg.AllowedHeaders = headerMap
	}
}

// WithAllowedHeadersAndValues sets the list of allowed headers in the CORSConfig.
// This function takes a map of string to string, representing the allowed headers and values.
// The map is merged with the existing AllowedHeaders field in the corsConfig struct.
func WithAllowedHeadersAndValues(headers map[string]string) CORSOption {
	return func(cfg *corsConfig) {
		for k, v := range headers {
			cfg.AllowedHeaders[k] = v
		}
	}
}

// WithPreflightHeaders sets the list of headers for preflight requests in the CORSConfig.
func WithPreflightHeaders(headers map[string]string) CORSOption {
	return func(cfg *corsConfig) {
		cfg.PreflightHeaders = headers
	}
}

// WithMaxAge sets the maximum age of a preflight request in seconds in the CORSConfig.
func WithMaxAge(maxAge int) CORSOption {
	return func(cfg *corsConfig) {
		cfg.MaxAge = maxAge
	}
}

/*
CORS is a middleware function that adds Cross-Origin Resource Sharing (CORS) headers to the HTTP response.

By default, it sets the Access-Control-Allow-Origin header to "*", which allows any origin to access the resource.
It also sets the Access-Control-Allow-Methods header to the HTTP methods defined in the Config, and the
Access-Control-Allow-Headers header to the HTTP headers defined in the Config.

The middleware can be customized by passing in one or more CORSOption values to the constructor. These options
can be used to configure the allowed origins, methods, headers, and other CORS settings.

Usage:

	// Create a new Router
	router := muxer.NewRouter()

	// Create a new CORS middleware with default options
	cors := muxer.CORS()

	// Register a new route with the CORS middleware
	router.HandleFunc("/api", myHandler).Methods("GET").Middleware(cors)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))

Alternatively, you can create a custom CORS middleware with specific options:

	// Create a new Router
	router := muxer.NewRouter()

	// Create a new CORS middleware with custom options
	cors := muxer.CORS(
		muxer.WithAllowedOrigins("https://example.com"),
		muxer.WithAllowedMethods("GET", "POST"),
		muxer.WithAllowedHeaders("Authorization", "Content-Type"),
		muxer.WithExposedHeaders("X-Custom-Header"),
		muxer.WithMaxAge(86400),
		muxer.WithAllowCredentials(),
	)

	// Register a new route with the custom CORS middleware
	router.HandleFunc("/api", myHandler).Methods("GET").Middleware(cors)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
*/
func CORS(options ...CORSOption) func(http.Handler) http.Handler {
	cfg := &corsConfig{
		AllowedHeaders:   make(map[string]string),
		PreflightHeaders: make(map[string]string),
	}

	for _, option := range options {
		option(cfg)
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && contains(cfg.AllowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// Always set the Access-Control-Allow-Origin header, even if the
				// incoming request does not contain an "Origin" header.
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			if len(cfg.AllowedMethods) > 0 {
				allowedMethods := strings.Join(cfg.AllowedMethods, ", ")
				w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			}

			if len(cfg.AllowedHeaders) > 0 {
				allowedHeaders := strings.Join(keys(cfg.AllowedHeaders), ", ")
				w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			}

			if r.Method == http.MethodOptions {
				if cfg.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(cfg.MaxAge), 10))
				}
				for k, v := range cfg.PreflightHeaders {
					w.Header().Set(k, v)
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// keys returns the keys of the given map as a string slice.
func keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// contains checks if the given string slice contains the given string.
func contains(slice []string, s string) bool {
	for _, elem := range slice {
		if elem == s {
			return true
		}
	}
	return false
}
