package muxer

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
This benchmark tests the performance of the Router implementation by simulating multiple requests to a single route.
The Router is first initialized and a single route is added to it using the AddRoute method. The route matches the
GET HTTP method and expects a parameter id in the URL path. When the route is matched, the id parameter is extracted
and written back to the response.

In the benchmark loop, a new request is created for each iteration, which matches the route and contains a different
id parameter value. The ServeHTTP method of the Router is called with the request and a new httptest.ResponseRecorder
to capture the response. The loop runs b.N times, which is a command-line flag that determines the number of times the
benchmark function should run.

The ResetTimer method is called before the benchmark loop to reset the timer and to exclude the time it takes to set
up the benchmark. The ServeHTTP method of the Router is called for each iteration and the time taken to process the
request is measured by the benchmarking framework.

The purpose of this benchmark is to measure the performance of the Router implementation under load and to identify
potential performance bottlenecks. By measuring the time taken to process a large number of requests, we can optimize
the implementation for better performance.
*/
func BenchmarkRouter(b *testing.B) {
	router := &Router{}

	// Add a route to the router
	router.HandleRoute(http.MethodGet, "/api/widgets/:widget/parts/:part/update", func(w http.ResponseWriter, r *http.Request) {
	})

	// Create a new request for the benchmark
	// req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	req := httptest.NewRequest(http.MethodGet, "/api/widgets/123/parts/456/update", nil)

	// Create a new recorder for the benchmark
	w := httptest.NewRecorder()

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, req)
	}
}

func BenchmarkManyPathVariables(b *testing.B) {
	router := &Router{}
	handler := func(w http.ResponseWriter, r *http.Request) {}
	router.HandleRoute(http.MethodGet, "/v1/:v1/:v2/:v3/:v4/:v5", handler)

	matchingRequest, _ := http.NewRequest("GET", "/v1/1/2/3/4/5", nil)
	recorder := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(recorder, matchingRequest)
	}
}
