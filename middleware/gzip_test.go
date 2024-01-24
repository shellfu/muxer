package middleware

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzip(t *testing.T) {
	tests := []struct {
		name                string
		contentType         string
		acceptEncoding      string
		expectedBody        string
		expectedContentType string
		expectedEncoding    string
	}{
		{
			name:                "gzip content encoding",
			contentType:         "text/plain",
			acceptEncoding:      "gzip",
			expectedBody:        "This is some sample text",
			expectedContentType: "text/plain",
			expectedEncoding:    "gzip",
		},
		{
			name:                "no gzip content encoding",
			contentType:         "text/plain",
			acceptEncoding:      "",
			expectedBody:        "This is some sample text",
			expectedContentType: "text/plain",
			expectedEncoding:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the given content type and body
			body := strings.NewReader(tc.expectedBody)
			req, err := http.NewRequest(http.MethodPost, "https://example.com", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tc.contentType)
			req.Header.Set("Accept-Encoding", tc.acceptEncoding)

			// Create a new recorder to capture the response
			rr := httptest.NewRecorder()

			// Create a handler with the Gzip middleware and a handler that echoes the request body
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", tc.contentType)
				if _, err := w.Write(body); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			})

			Gzip(handler).ServeHTTP(rr, req)

			// Check that the response has the expected content type and encoding
			if rr.Header().Get("Content-Type") != tc.expectedContentType {
				t.Errorf("expected Content-Type %q, got %q", tc.expectedContentType, rr.Header().Get("Content-Type"))
			}
			if rr.Header().Get("Content-Encoding") != tc.expectedEncoding {
				t.Errorf("expected Content-Encoding %q, got %q", tc.expectedEncoding, rr.Header().Get("Content-Encoding"))
			}

			// Decode the response body if it's gzip encoded
			var bodyBytes []byte
			if rr.Header().Get("Content-Encoding") == "gzip" {
				reader, err := gzip.NewReader(rr.Body)
				if err != nil {
					t.Fatal(err)
				}
				defer reader.Close()
				bodyBytes, err = ioutil.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				bodyBytes = rr.Body.Bytes()
			}

			// Check that the response body is the same as the expected body
			if !bytes.Equal(bodyBytes, []byte(tc.expectedBody)) {
				t.Errorf("expected body %q, got %q", tc.expectedBody, string(bodyBytes))
			}
		})
	}
}
