# Implementing Prometheus Metrics in Go HTTP Server

Prometheus is an open-source monitoring system with a dimensional data model,
flexible query language, efficient time series database and modern alerting
approach. This tutorial will guide you on how to implement Prometheus metrics
in a Go HTTP server using the `muxer` router.

## Prerequisites

Before you start, you need to have Go installed in your environment. Also, make
sure to have your `muxer` router defined. 

## Step 1 — Installing Prometheus Client Library

To start, you need to install the Prometheus client library for Go. You can do
it with the following `go get` command:

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
```

## Step 2 — Defining Prometheus Metrics

In a new file or where your custom middleware is defined, you need to import 
the Prometheus client library and define the metrics that you want to observe. 
For instance, we will create two metrics: one for observing HTTP requests' 
duration and another for counting the number of received requests.

```go
import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    mHTTPDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests.",
        },
        []string{"service", "path"},
    )

    mHTTPCount = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_request_count",
            Help: "Number of requests received.",
        },
        []string{"service", "method"},
    )
)
```

## Step 3 — Implementing Prometheus Middleware

Next, we will create two middleware functions to track these metrics. These
functions will wrap the HTTP handlers to measure and record the request
duration and count.

```go
func withHTTPDuration(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            path := r.URL.Path
            prometheus.NewTimer(mHTTPDuration.WithLabelValues("my-service", path)).ObserveDuration()
        }()
        next.ServeHTTP(w, r)
    })
}

func withHTTPCount(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer mHTTPCount.WithLabelValues("my-service", r.Method).Inc()
        next.ServeHTTP(w, r)
    })
}
```

Please replace "my-service" with the actual service name.

## Step 4 — Using Prometheus Middleware in Router

With the `muxer` package, we can add middleware functions to the router using
the `Use` method. We will add our Prometheus middleware to the router so that
it tracks metrics for every request.

```go
router := muxer.NewRouter()
router.Use(withHTTPDuration, withHTTPCount)
```

## Step 5 — Exposing Prometheus Metrics

Lastly, we need to expose these metrics so that Prometheus can scrape them. We
can do it by adding a new route to our server that serves the metrics.

```go
router.Handle(http.MethodGet, "/metrics", promhttp.Handler())
```

`muxer` is now ready to track and expose Prometheus metrics!

## Conclusion

In this guide, we showed you how to implement Prometheus metrics using `muxer`.
This setup will provide you with basic insights about your service, such as 
request durations and request counts. By analyzing these metrics, you
can better understand your service performance.
