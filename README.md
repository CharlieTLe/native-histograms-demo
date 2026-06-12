# Native Histogram Demo

A minimal demo showing Prometheus scraping native histograms from a Go application, orchestrated with docker-compose.

## What's Inside

- **demo-app** — A Go HTTP server that registers a native histogram (`demo_request_duration_seconds`) with `NativeHistogramBucketFactor: 1.1` and observes random exponentially-distributed latencies every 100ms.
- **prometheus** — A Prometheus instance configured with `scrape_native_histograms: true` to scrape the app using protobuf content negotiation.

## Quick Start

```bash
docker compose up --build
```

- App metrics: http://localhost:8080/metrics
- Prometheus UI: http://localhost:9090

## Example Queries

Once Prometheus has scraped a few samples, try these in the Prometheus UI:

```promql
# Raw native histogram
demo_request_duration_seconds

# 95th percentile latency
histogram_quantile(0.95, demo_request_duration_seconds)

# Average observed value
histogram_avg(demo_request_duration_seconds)

# Total observation count
histogram_count(demo_request_duration_seconds)

# Fraction of observations between 0 and 500ms
histogram_fraction(0, 0.5, demo_request_duration_seconds)
```

## Project Structure

```
├── docker-compose.yml
├── prometheus.yml
└── app/
    ├── Dockerfile
    ├── go.mod
    └── main.go
```
