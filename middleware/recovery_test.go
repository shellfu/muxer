package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockLogger struct {
	buf bytes.Buffer
}

func (l *mockLogger) Println(v ...interface{}) {
	for _, msg := range v {
		l.buf.WriteString(msg.(string))
	}
}

func TestRecoveryHandler(t *testing.T) {
	logger := &mockLogger{}
	// Define test cases
	tests := []struct {
		name       string
		printStack bool
		handler    http.Handler
		logger     RecoveryLogger
	}{
		{
			name:       "panic with stack trace",
			printStack: true,
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected error")
			}),
			logger: logger,
		},
		{
			name:       "panic without stack trace",
			printStack: false,
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected error")
			}),
			logger: logger,
		},
		{
			name:       "no panic",
			printStack: false,
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			logger: logger,
		},
		{
			name:       "no logger without stack trace",
			printStack: false,
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected error")
			}),
			logger: nil,
		},
		{
			name:       "no logger with stack trace",
			printStack: true,
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected error")
			}),
			logger: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the recovery middleware
			handler := RecoveryHandler(tt.logger, tt.printStack)(tt.handler)

			// Create a request with the middleware
			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// Check the response status code
			if rec.Code != http.StatusInternalServerError && tt.name != "no panic" {
				t.Errorf("unexpected status code: %v", rec.Code)
			}

			// Check if the error message and stack trace are logged
			if tt.printStack {
				if logger.buf.String() == "" {
					t.Errorf("error message and stack trace not logged")
				}
			}
		})
	}
}
