package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	requestDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:                            "demo_request_duration_seconds",
		Help:                            "Simulated request duration in seconds (native histogram).",
		NativeHistogramBucketFactor:      1.1,
		NativeHistogramMaxBucketNumber:   100,
		NativeHistogramMinResetDuration:  1 * time.Hour,
		NativeHistogramZeroThreshold:     0.001,
	})
	prometheus.MustRegister(requestDuration)

	go func() {
		for {
			// Exponentially-distributed latency with mean ~100ms, range roughly 10ms–2s.
			latency := rand.ExpFloat64() * 0.1
			requestDuration.Observe(latency)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
