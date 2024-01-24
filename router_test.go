package muxer

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	. "github.com/shellfu/muxer/middleware"
)

func TestRouter(t *testing.T) {
	router := NewRouter()

	testCases := []struct {
		method       string
		path         string
		expectedCode int
		expectedBody string
	}{
		{http.MethodGet, "/", http.StatusNotFound, "404 page not found"},
		{http.MethodGet, "/users", http.StatusNotFound, "404 page not found"},
		{http.MethodGet, "/users/123", http.StatusOK, "123"},
		{http.MethodGet, "/users/123.js", http.StatusOK, "123.js"},
		{http.MethodGet, "/users/123-js", http.StatusOK, "123-js"},
		{http.MethodGet, "/users/123_js", http.StatusOK, "123_js"},
		{http.MethodGet, "/users/abc", http.StatusOK, "abc"},
		{http.MethodPost, "/users/123", http.StatusMethodNotAllowed, "Method not allowed"},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		w := httptest.NewRecorder()

		router.HandleRoute(http.MethodGet, "/users/:id", func(w http.ResponseWriter, r *http.Request) {
			id := router.Params(r)["id"]
			if _, err := w.Write([]byte(id)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		router.ServeHTTP(w, req)

		if w.Code != tc.expectedCode {
			t.Errorf("unexpected status code: expected=%d, actual=%d", tc.expectedCode, w.Code)
		}

		if strings.Replace(w.Body.String(), "\n", "", -1) != tc.expectedBody {
			t.Errorf("unexpected response body: expected=%s, actual=%s", tc.expectedBody, w.Body.String())
		}
	}
}

func TestRouter_Handle(t *testing.T) {
	router := NewRouter()

	testCases := []struct {
		method  string
		path    string
		handler http.Handler
	}{
		{http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})},
		{http.MethodPost, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})},
		{http.MethodDelete, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})},
	}

	for _, tc := range testCases {
		router.Handle(tc.method, tc.path, tc.handler)
	}

	if len(router.routes) != len(testCases) {
		t.Errorf("unexpected number of routes: expected=%d, actual=%d", len(testCases), len(router.routes))
	}

	for i, route := range router.routes {
		tc := testCases[i]

		if route.method != tc.method {
			t.Errorf("unexpected method for route %d: expected=%s, actual=%s", i, tc.method, route.method)
		}

		if route.path.String() != "^"+tc.path+"$" {
			t.Errorf("unexpected path for route %d: expected=%s, actual=%s", i, "^"+tc.path+"$", route.path.String())
		}

		if route.handler == nil {
			t.Errorf("unexpected nil handler for route %d", i)
		}
	}
}

func TestParams(t *testing.T) {
	testCases := []struct {
		path           string
		params         map[string]string
		expectedParams map[string]string
	}{
		{
			path:           "/users/123",
			params:         map[string]string{"id": "123"},
			expectedParams: map[string]string{"id": "123"},
		},
		{
			path:           "/users/456",
			params:         map[string]string{},
			expectedParams: map[string]string{},
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)

		// Set params in context
		if len(tc.params) > 0 {
			ctx := context.WithValue(req.Context(), paramsKey, tc.params)
			req = req.WithContext(ctx)
		}

		router := NewRouter()
		actualParams := router.Params(req)

		if len(actualParams) != len(tc.expectedParams) {
			t.Errorf("unexpected number of params for path %s: expected=%d, actual=%d", tc.path, len(tc.expectedParams), len(actualParams))
		}

		for k, v := range tc.expectedParams {
			if actualParams[k] != v {
				t.Errorf("unexpected param value for path %s: key=%s, expected=%s, actual=%s", tc.path, k, v, actualParams[k])
			}
		}

		if len(actualParams) == 0 && len(tc.params) > 0 {
			t.Errorf("params not found in context for path %s", tc.path)
		}
	}
}

func TestSubrouter(t *testing.T) {
	router := NewRouter()

	// Create subrouter for www.example.com
	example := router.Subrouter("www.example.com")
	example.HandlerFunc(http.MethodGet, "/example", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Example")
	})

	// Create subrouter for /api
	api := router.Subrouter("/api")
	api.HandleRoute(http.MethodGet, "/users/:id", func(w http.ResponseWriter, req *http.Request) {
		params := api.Params(req)
		id := params["id"]
		if id != "123" {
			t.Errorf("Unexpected ID: %s", id)
		}
	})

	// Define test cases
	tests := []struct {
		method string
		path   string
		status int
		isHost bool
	}{
		{http.MethodGet, "/example", http.StatusOK, true},
		{http.MethodGet, "/api/users/123", http.StatusOK, false},
		{http.MethodGet, "/users/123", http.StatusNotFound, false},
	}

	// Run test cases
	for _, tc := range tests {
		req, err := http.NewRequest(tc.method, tc.path, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Set host value for www.example.com subrouter
		if tc.isHost {
			req.URL.Host = "www.example.com"
		}

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != tc.status {
			t.Errorf("Unexpected status: %d", recorder.Code)
		}
	}
}

func TestRouter_HandleRoute(t *testing.T) {
	router := NewRouter()

	testCases := []struct {
		method       string
		path         string
		handlerFunc  http.HandlerFunc
		expectedCode int
	}{
		{http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK},
		{http.MethodPost, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK},
		{http.MethodDelete, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK},
	}

	for _, tc := range testCases {
		router.HandleRoute(tc.method, tc.path, tc.handlerFunc)
	}

	if len(router.routes) != len(testCases) {
		t.Errorf("unexpected number of routes: expected=%d, actual=%d", len(testCases), len(router.routes))
	}
}

func TestRouter_HandlerFuncWithMethods(t *testing.T) {
	router := NewRouter()

	testCases := []struct {
		method       string
		path         string
		handlerFunc  http.HandlerFunc
		expectedCode int
	}{
		{http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK},
		{http.MethodPost, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK},
		{http.MethodDelete, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK},
	}

	for _, tc := range testCases {
		router.HandlerFuncWithMethods([]string{tc.method}, tc.path, tc.handlerFunc)
	}

	if len(router.routes) != len(testCases) {
		t.Errorf("unexpected number of routes: expected=%d, actual=%d", len(testCases), len(router.routes))
	}
}

func TestRouter_Use(t *testing.T) {
	router := NewRouter()

	testMiddleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-Test-1", "1")
			next.ServeHTTP(w, r)
		})
	}

	testMiddleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-Test-2", "2")
			next.ServeHTTP(w, r)
		})
	}

	router.Use(testMiddleware1, testMiddleware2)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	router.HandleRoute(http.MethodGet, "/", handlerFunc)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(rr, req)

	if rr.Header().Get("X-Middleware-Test-1") != "1" || rr.Header().Get("X-Middleware-Test-2") != "2" {
		t.Errorf("unexpected middleware headers: expected=%s, actual=%s", "1, 2", rr.Header().Get("X-Middleware-Test"))
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	router := NewRouter()

	testCases := []struct {
		method       string
		path         string
		handlerFunc  http.HandlerFunc
		expectedCode int
		expectedBody string
	}{
		{http.MethodGet, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := router.Params(r)["id"]
			if _, err := w.Write([]byte(id)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}), http.StatusOK, "123"},
		{http.MethodPost, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK, ""},
		{http.MethodDelete, "/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), http.StatusOK, ""},
	}

	for _, tc := range testCases {
		router.HandleRoute(tc.method, tc.path, tc.handlerFunc)

		req := httptest.NewRequest(tc.method, "/users/123", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != tc.expectedCode {
			t.Errorf("unexpected status code: expected=%d, actual=%d", tc.expectedCode, w.Code)
		}

		if strings.Replace(w.Body.String(), "\n", "", -1) != tc.expectedBody {
			t.Errorf("unexpected response body: expected=%s, actual=%s", tc.expectedBody, w.Body.String())
		}
	}
}

func TestNotFoundHandler(t *testing.T) {
	notFoundHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Custom 404 Page")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router := NewRouter(WithNotFoundHandler(notFoundHandlerFunc))
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := router.Params(r)["id"]
		if _, err := w.Write([]byte(id)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
	router.HandleRoute(http.MethodGet, "/users/:id", handlerFunc)

	testCases := []struct {
		path         string
		expectedCode int
		expectedBody string
	}{
		{"/non-existing-path", http.StatusNotFound, "Custom 404 Page"},
		{"/users/123", http.StatusOK, "123"},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != tc.expectedCode {
			t.Errorf("Expected status code: %d. Got: %d", tc.expectedCode, resp.Code)
		}
		if resp.Body.String() != tc.expectedBody {
			t.Errorf("Expected response body: %s. Got: %s", tc.expectedBody, resp.Body.String())
		}
	}
}

func TestMaxRequestBodySize(t *testing.T) {
	maxRequestBodySize := int64(1024)
	router := NewRouter(WithMaxRequestBodySize(maxRequestBodySize))

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.HandleRoute(http.MethodPost, "/users/:id", handlerFunc)

	testCases := []struct {
		path         string
		body         io.Reader
		expectedCode int
	}{
		{"/users/123", strings.NewReader(strings.Repeat("a", int(maxRequestBodySize+1))), http.StatusRequestEntityTooLarge},
		{"/users/123", strings.NewReader(strings.Repeat("a", int(maxRequestBodySize))), http.StatusOK},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodPost, tc.path, tc.body)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		if resp.Code != tc.expectedCode {
			t.Errorf("Expected status code: %d. Got: %d", tc.expectedCode, resp.Code)
		}
	}
}

func TestHandlerFunc(t *testing.T) {
	router := NewRouter()

	// Test adding a route with HandlerFunc
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello, world!")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	router.HandlerFunc("GET", "/hello", handlerFunc)

	// Test that the route works
	req, err := http.NewRequest("GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("HandlerFunc route returned wrong status code: got %v, want %v", rr.Code, http.StatusOK)
	}
	if body := rr.Body.String(); body != "Hello, world!" {
		t.Errorf("HandlerFunc route returned unexpected body: got %v, want %v", body, "Hello, world!")
	}
}

func TestEnableCORSOption(t *testing.T) {
	tests := []struct {
		name             string
		origin           string
		expectedHeaders  map[string][]string
		expectedMaxAge   string
		enableCORSOption []CORSOption
	}{
		{
			name:   "CORS headers set correctly",
			origin: "http://example.com",
			expectedHeaders: map[string][]string{
				"Access-Control-Allow-Origin":  {"http://example.com"},
				"Access-Control-Allow-Headers": {"Content-Type"},
			},
			enableCORSOption: []CORSOption{
				WithAllowedOrigins("http://example.com"),
				WithAllowedHeaders("Content-Type"),
			},
		},
		{
			name:            "CORS headers not set if no origin",
			expectedHeaders: map[string][]string{},
			enableCORSOption: []CORSOption{
				WithAllowedOrigins("http://example.com"),
				WithAllowedHeaders("Content-Type"),
			},
		},
		{
			name:             "CORS headers not set if origin not allowed",
			origin:           "http://example2.com",
			expectedHeaders:  map[string][]string{},
			enableCORSOption: []CORSOption{WithAllowedOrigins("http://example.com")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			router := NewRouter()
			router.Use(CORS(tc.enableCORSOption...))

			router.HandlerFunc(http.MethodGet, "/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if _, err := w.Write([]byte(`{"message": "hello world"}`)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}))

			req, err := http.NewRequest(http.MethodGet, "http://example.com/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}

			router.ServeHTTP(rr, req)

			// Check headers
			actualHeaders := rr.Header()
			for k, v := range tc.expectedHeaders {
				actual := actualHeaders[k]
				if !reflect.DeepEqual(actual, v) {
					t.Errorf("expected header %s with value %v, got %v", k, v, actual)
				}
			}
		})
	}
}
