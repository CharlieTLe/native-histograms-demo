# Native Histogram Demo

A minimal demo showing Prometheus scraping native histograms from a Go application, orchestrated with docker-compose.

## What's Inside

- **demo-app** — A Go HTTP server that registers both a native histogram (`demo_request_duration_seconds`) and a classic histogram (`demo_classic_request_duration_seconds`) observing the same random exponentially-distributed latencies every 100ms, enabling side-by-side comparison.
- **prometheus** — A Prometheus instance configured with `scrape_native_histograms: true` to scrape the app using protobuf content negotiation.

## Quick Start

```bash
docker compose up --build
```

- App metrics: http://localhost:8080/metrics
- Prometheus UI: http://localhost:9090

## Example Queries

Once Prometheus has scraped a few samples, try these in the Prometheus UI (click to open with query pre-filled):

### Native Histogram

- [Raw native histogram](http://localhost:9090/query?g0.expr=demo_request_duration_seconds&g0.show_tree=0&g0.tab=table&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  demo_request_duration_seconds
  ```
- [95th percentile latency](http://localhost:9090/query?g0.expr=histogram_quantile%280.95%2C+demo_request_duration_seconds%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_quantile(0.95, demo_request_duration_seconds)
  ```
- [Average observed value](http://localhost:9090/query?g0.expr=histogram_avg%28demo_request_duration_seconds%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_avg(demo_request_duration_seconds)
  ```
- [Total observation count](http://localhost:9090/query?g0.expr=histogram_count%28demo_request_duration_seconds%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_count(demo_request_duration_seconds)
  ```
- [Fraction of observations between 0 and 500ms](http://localhost:9090/query?g0.expr=histogram_fraction%280%2C+0.5%2C+demo_request_duration_seconds%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_fraction(0, 0.5, demo_request_duration_seconds)
  ```

### Classic Histogram

- [95th percentile latency (classic)](http://localhost:9090/query?g0.expr=histogram_quantile%280.95%2C+rate%28demo_classic_request_duration_seconds_bucket%5B5m%5D%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_quantile(0.95, rate(demo_classic_request_duration_seconds_bucket[5m]))
  ```
- [Average observed value (classic)](http://localhost:9090/query?g0.expr=rate%28demo_classic_request_duration_seconds_sum%5B5m%5D%29+%2F+rate%28demo_classic_request_duration_seconds_count%5B5m%5D%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  rate(demo_classic_request_duration_seconds_sum[5m]) / rate(demo_classic_request_duration_seconds_count[5m])
  ```
- [Total observation count (classic)](http://localhost:9090/query?g0.expr=demo_classic_request_duration_seconds_count&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  demo_classic_request_duration_seconds_count
  ```

### Side-by-Side Comparison

Compare the 95th percentile accuracy between native and classic histograms — [open both panels](http://localhost:9090/query?g0.expr=histogram_quantile%280.95%2C+demo_request_duration_seconds%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0&g1.expr=histogram_quantile%280.95%2C+rate%28demo_classic_request_duration_seconds_bucket%5B5m%5D%29%29&g1.show_tree=0&g1.tab=graph&g1.range_input=1h&g1.res_type=auto&g1.res_density=medium&g1.display_mode=lines&g1.show_exemplars=0):

```promql
# Native — higher resolution, no bucket boundary errors
histogram_quantile(0.95, demo_request_duration_seconds)

# Classic — limited to default bucket boundaries (.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)
histogram_quantile(0.95, rate(demo_classic_request_duration_seconds_bucket[5m]))
```

The demo app generates exponentially-distributed latencies with mean 100ms, giving a theoretical 95th percentile of **299.6ms**[^1]. In practice:

| Histogram | 95th percentile | Notes |
|-----------|----------------|-------|
| Native    | ~283ms | Close to theoretical value due to fine-grained exponential bucketing |
| Classic   | ~340ms | Overshoots because it interpolates linearly between the fixed 0.25s and 0.5s bucket boundaries |

## Project Structure

```
├── docker-compose.yml
├── prometheus.yml
└── app/
    ├── Dockerfile
    ├── go.mod
    └── main.go
```

## Cleanup

```bash
docker compose down -v
```

[^1]: For an exponential distribution with mean μ, the CDF is `F(x) = 1 − e^(−x/μ)`. Solving `F(x) = 0.95` for μ = 0.1s: `x = −0.1 × ln(0.05) ≈ 0.2996s`.
