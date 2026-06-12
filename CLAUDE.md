# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A demo comparing Prometheus native histograms vs classic histograms, scraped from a Go application. Two services run via Docker Compose: a Go demo-app (port 8080) that generates synthetic request latency data with both histogram types, and Prometheus v3 (port 9090) that scrapes those metrics (including itself) every 1s.

## Commands

### Run the demo
```
docker compose up --build
```

### Build the Go app locally
```
cd app && CGO_ENABLED=0 go build -o demo-app .
```

### Cleanup
```
docker compose down -v
```

### Access points
- Metrics endpoint: http://localhost:8080/metrics
- Prometheus UI: http://localhost:9090

### Query the Prometheus API
```
curl -s --get 'http://localhost:9090/api/v1/query' --data-urlencode 'query=<promql>'
```

## Architecture

The Go app (`app/main.go`) creates two histogram vectors with labels `method`, `handler`, `status`, `service`, `region` (432 combinations):
- `demo_request_duration_seconds` — native histogram (exponential bucketing, factor 1.1, max 100 buckets)
- `demo_classic_request_duration_seconds` — classic histogram (`prometheus.DefBuckets`)

Both observe the same exponentially-distributed latencies every 100ms, with varying means per handler (`/health` 10ms, `/api/users` 100ms, `/api/products` 120ms, `/api/orders` 150ms).

Prometheus is configured in `prometheus.yml` with `scrape_native_histograms: true` and a 1s scrape interval. It also scrapes itself at `localhost:9090` for TSDB storage analysis.

## Dependencies

- Go 1.25 with `github.com/prometheus/client_golang` v1.23.2
- Docker multi-stage build (Go 1.25-alpine builder, alpine:3.20 runtime)
