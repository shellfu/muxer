package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		slice  []string
		item   string
		result bool
	}{
		{
			name:   "Item is present in slice",
			slice:  []string{"foo", "bar", "baz"},
			item:   "bar",
			result: true,
		},
		{
			name:   "Item is not present in slice",
			slice:  []string{"foo", "bar", "baz"},
			item:   "qux",
			result: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if contains(tc.slice, tc.item) != tc.result {
				t.Errorf("contains(%v, %q) = %v, want %v", tc.slice, tc.item, !tc.result, tc.result)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		origin             string
		config             []CORSOption
		expectedStatusCode int
		expectedHeaders    http.Header
	}{
		{
			name:   "Request without Origin header",
			method: http.MethodGet,
			origin: "",
			config: []CORSOption{
				WithAllowedMethods(http.MethodGet, http.MethodPost),
			},
			expectedStatusCode: http.StatusOK,
			expectedHeaders: http.Header{
				"Access-Control-Allow-Methods": []string{http.MethodGet, http.MethodPost},
			},
		},
		{
			name:   "Simple request with matching Origin header",
			method: http.MethodGet,
			origin: "http://example.com",
			config: []CORSOption{
				WithAllowedOrigins("http://example.com"),
				WithAllowedMethods(http.MethodGet, http.MethodPost),
			},
			expectedStatusCode: http.StatusOK,
			expectedHeaders: http.Header{
				"Access-Control-Allow-Origin":  []string{"http://example.com"},
				"Access-Control-Allow-Methods": []string{http.MethodGet, http.MethodPost},
			},
		},
		{
			name:   "Preflight request with matching Origin header",
			method: http.MethodOptions,
			origin: "http://example.com",
			config: []CORSOption{
				WithAllowedOrigins("http://example.com"),
				WithAllowedMethods(http.MethodGet, http.MethodPost),
				WithAllowedHeaders("X-Custom-Header-1"),
				WithAllowedHeadersAndValues(map[string]string{
					"X-Custom-Header": "foobar",
				}),
				WithPreflightHeaders(map[string]string{
					"X-Preflight-Header": "123",
				}),
				WithMaxAge(3600),
			},
			expectedStatusCode: http.StatusOK,
			expectedHeaders: http.Header{
				"Access-Control-Allow-Origin":  []string{"http://example.com"},
				"Access-Control-Allow-Methods": []string{http.MethodGet, http.MethodPost},
				"Access-Control-Allow-Headers": []string{"X-Custom-Header-1", "X-Custom-Header"},
				"Access-Control-Max-Age":       []string{"3600"},
				"X-Preflight-Header":           []string{"123"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, "http://example.com", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}

			rr := httptest.NewRecorder()
			handler := CORS(tc.config...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Custom-Header", "foobar")
			}))
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, rr.Code)
			}

			// Check headers
			for k, v := range tc.expectedHeaders {
				expectedValue := strings.Join(v, ", ")
				if got := rr.Header().Get(k); got != expectedValue {
					t.Errorf("expected header %s with value '%s', got '%s'", k, expectedValue, got)
				}
			}
		})
	}
}
