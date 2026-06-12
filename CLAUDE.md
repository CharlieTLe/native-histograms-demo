# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A minimal demo showing Prometheus native histograms scraped from a Go application. Two services run via Docker Compose: a Go demo-app (port 8080) that generates synthetic request latency data with native histogram metrics, and Prometheus v3 (port 9090) that scrapes those metrics using protobuf content negotiation.

## Commands

### Run the demo
```
docker compose up --build
```

### Build the Go app locally
```
cd app && CGO_ENABLED=0 go build -o demo-app .
```

### Access points
- Metrics endpoint: http://localhost:8080/metrics
- Prometheus UI: http://localhost:9090

## Architecture

The Go app (`app/main.go`) creates a single native histogram metric `demo_request_duration_seconds` configured with exponential bucketing (factor 1.1, max 100 buckets). A goroutine generates exponentially-distributed latencies (~10ms-2s, mean ~100ms) and observes them every 100ms. Metrics are exposed via `promhttp.Handler()`.

Prometheus is configured in `prometheus.yml` with `scrape_native_histograms: true` to enable native histogram ingestion, scraping the demo-app every 15s.

## Dependencies

- Go 1.25 with `github.com/prometheus/client_golang` v1.23.2
- Docker multi-stage build (Go 1.25-alpine builder, alpine:3.20 runtime)
