package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	methods  = []string{"GET", "POST", "PUT", "DELETE"}
	handlers = []string{"/api/users", "/api/orders", "/api/products", "/health"}
	statuses = []string{"200", "400", "500"}
	services = []string{"web", "api", "worker"}
	regions  = []string{"us-east", "us-west", "eu-west"}

	// Mean latency multiplier per handler to simulate realistic variance.
	handlerLatency = map[string]float64{
		"/api/users":    0.1,  // 100ms mean
		"/api/orders":   0.15, // 150ms mean
		"/api/products": 0.12, // 120ms mean
		"/health":       0.01, // 10ms mean
	}
)

func main() {
	labels := []string{"method", "handler", "status", "service", "region"}

	nativeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            "demo_request_duration_seconds",
		Help:                            "Simulated request duration in seconds (native histogram).",
		NativeHistogramBucketFactor:     1.1,
		NativeHistogramMaxBucketNumber:  100,
		NativeHistogramMinResetDuration: 1 * time.Hour,
		NativeHistogramZeroThreshold:    0.001,
	}, labels)

	classicHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "demo_classic_request_duration_seconds",
		Help:    "Simulated request duration in seconds (classic histogram).",
		Buckets: prometheus.DefBuckets,
	}, labels)

	prometheus.MustRegister(nativeHistogram, classicHistogram)

	go func() {
		for {
			method := methods[rand.Intn(len(methods))]
			handler := handlers[rand.Intn(len(handlers))]
			status := statuses[rand.Intn(len(statuses))]
			service := services[rand.Intn(len(services))]
			region := regions[rand.Intn(len(regions))]

			latency := rand.ExpFloat64() * handlerLatency[handler]

			nativeHistogram.WithLabelValues(method, handler, status, service, region).Observe(latency)
			classicHistogram.WithLabelValues(method, handler, status, service, region).Observe(latency)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
