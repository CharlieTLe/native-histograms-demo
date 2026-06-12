# Native Histogram Demo

A minimal demo showing Prometheus scraping native histograms from a Go application, orchestrated with docker-compose.

## What's Inside

- **demo-app** — A Go HTTP server that registers both a native histogram (`demo_request_duration_seconds`) and a classic histogram (`demo_classic_request_duration_seconds`) with labels `method`, `handler`, `status`, `service`, and `region`, observing the same random exponentially-distributed latencies every 100ms across 432 label combinations.
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

- [Raw native histogram (single series)](http://localhost:9090/query?g0.expr=demo_request_duration_seconds%7Bhandler%3D%22%2Fapi%2Fusers%22%2Cmethod%3D%22GET%22%2Cservice%3D%22web%22%2Cregion%3D%22us-east%22%2Cstatus%3D%22200%22%7D&g0.show_tree=0&g0.tab=table&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  demo_request_duration_seconds{handler="/api/users",method="GET",service="web",region="us-east",status="200"}
  ```
- [95th percentile latency by handler](http://localhost:9090/query?g0.expr=histogram_quantile%280.95%2C+sum+by+%28handler%29+%28demo_request_duration_seconds%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_quantile(0.95, sum by (handler) (demo_request_duration_seconds))
  ```
- [Average observed value by handler](http://localhost:9090/query?g0.expr=histogram_avg%28sum+by+%28handler%29+%28demo_request_duration_seconds%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_avg(sum by (handler) (demo_request_duration_seconds))
  ```
- [Total observation count by handler](http://localhost:9090/query?g0.expr=histogram_count%28sum+by+%28handler%29+%28demo_request_duration_seconds%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_count(sum by (handler) (demo_request_duration_seconds))
  ```
- [Fraction of observations between 0 and 500ms by handler](http://localhost:9090/query?g0.expr=histogram_fraction%280%2C+0.5%2C+sum+by+%28handler%29+%28demo_request_duration_seconds%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_fraction(0, 0.5, sum by (handler) (demo_request_duration_seconds))
  ```

### Classic Histogram

- [95th percentile latency by handler (classic)](http://localhost:9090/query?g0.expr=histogram_quantile%280.95%2C+sum+by+%28handler%2C+le%29+%28rate%28demo_classic_request_duration_seconds_bucket%5B5m%5D%29%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  histogram_quantile(0.95, sum by (handler, le) (rate(demo_classic_request_duration_seconds_bucket[5m])))
  ```
- [Average observed value by handler (classic)](http://localhost:9090/query?g0.expr=sum+by+%28handler%29+%28rate%28demo_classic_request_duration_seconds_sum%5B5m%5D%29%29+%2F+sum+by+%28handler%29+%28rate%28demo_classic_request_duration_seconds_count%5B5m%5D%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  sum by (handler) (rate(demo_classic_request_duration_seconds_sum[5m])) / sum by (handler) (rate(demo_classic_request_duration_seconds_count[5m]))
  ```
- [Total observation count by handler (classic)](http://localhost:9090/query?g0.expr=sum+by+%28handler%29+%28demo_classic_request_duration_seconds_count%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0)
  ```promql
  sum by (handler) (demo_classic_request_duration_seconds_count)
  ```

### Side-by-Side Comparison

Compare the 95th percentile accuracy between native and classic histograms — [open both panels](http://localhost:9090/query?g0.expr=histogram_quantile%280.95%2C+sum+by+%28handler%29+%28demo_request_duration_seconds%29%29&g0.show_tree=0&g0.tab=graph&g0.range_input=1h&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0&g1.expr=histogram_quantile%280.95%2C+sum+by+%28handler%2C+le%29+%28rate%28demo_classic_request_duration_seconds_bucket%5B5m%5D%29%29%29&g1.show_tree=0&g1.tab=graph&g1.range_input=1h&g1.res_type=auto&g1.res_density=medium&g1.display_mode=lines&g1.show_exemplars=0):

```promql
# Native — higher resolution, no bucket boundary errors
histogram_quantile(0.95, sum by (handler) (demo_request_duration_seconds))

# Classic — limited to default bucket boundaries (.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10)
histogram_quantile(0.95, sum by (handler, le) (rate(demo_classic_request_duration_seconds_bucket[5m])))
```

The demo app generates exponentially-distributed latencies with varying means per handler (`/health` 10ms, `/api/users` 100ms, `/api/products` 120ms, `/api/orders` 150ms), giving different theoretical 95th percentiles[^1]. In practice:

| Histogram | 95th percentile | Notes |
|-----------|----------------|-------|
| Native    | ~285ms (`/api/users`) | Close to theoretical value due to fine-grained exponential bucketing |
| Classic   | ~340ms (`/api/users`) | Overshoots because it interpolates linearly between the fixed 0.25s and 0.5s bucket boundaries |

### Cardinality

With 5 labels (`method`, `handler`, `status`, `service`, `region`) producing 432 combinations (4 × 4 × 3 × 3 × 3), the classic histogram creates **14 time series per combination** (11 `_bucket` + `_sum` + `_count` + `_created`), while the native histogram uses **1 series per combination**:

| | Classic | Native |
|---|---|---|
| Series per combination | 14 | 1 |
| Total time series | ~6,048 | ~432 |

The classic histogram produces **14x more time series** for the same data with fewer buckets (11 fixed vs ~80 exponential).

### Query Performance

With 432 label combinations, the 95th percentile query is **~2x slower** for classic histograms[^6]:

| | Native | Classic |
|---|---|---|
| Average | ~24ms | ~49ms |

The classic query must fetch and `rate()` across 14 series per label combination before computing the quantile, while the native query operates on a single series per combination.

### Storage

The classic histogram dominates TSDB resource usage[^7]:

| | Classic | Native |
|---|---|---|
| Series | 6,048 (82% of TSDB) | 432 (6% of TSDB) |
| Samples per scrape | 6,048 float | 432 histogram |
| `le` label index memory | 44 KB | 0 (not needed) |

Each scrape ingests **14x more samples** for the classic histogram. The `le` label — which only exists for classic `_bucket` series — adds 44 KB of label index overhead[^8]. Native histograms avoid this entirely by encoding bucket data within each sample using protobuf.

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

## Converting a Classic Histogram to Native

### 1. Update the histogram options

Replace the fixed `Buckets` with native histogram fields:

```go
// Before (classic)
prometheus.NewHistogram(prometheus.HistogramOpts{
    Name:    "request_duration_seconds",
    Help:    "Request duration in seconds.",
    Buckets: prometheus.DefBuckets,
})

// After (native)
prometheus.NewHistogram(prometheus.HistogramOpts{
    Name:                           "request_duration_seconds",
    Help:                           "Request duration in seconds.",
    NativeHistogramBucketFactor:    1.1,
    NativeHistogramMaxBucketNumber: 100,
    NativeHistogramMinResetDuration: 1 * time.Hour,
    NativeHistogramZeroThreshold:   0.001,
})
```

- `NativeHistogramBucketFactor`[^2] — controls bucket width. `1.1` means each bucket boundary is 1.1x the previous. Lower = more precision, more buckets.
- `NativeHistogramMaxBucketNumber`[^3] — caps the bucket count. If exceeded, buckets get merged (resolution decreases).
- `NativeHistogramMinResetDuration`[^4] — minimum time before the histogram resets after a bucket merge.
- `NativeHistogramZeroThreshold`[^5] — observations below this value go into the special zero bucket.

You can keep `Buckets` alongside the native fields to emit both formats during a transition period. Remove `Buckets` once you're fully migrated.

### 2. Enable native histogram scraping in Prometheus

Add `scrape_native_histograms: true` to your scrape config:

```yaml
scrape_configs:
  - job_name: my-app
    scrape_native_histograms: true
    static_configs:
      - targets: ['my-app:8080']
```

Without this, Prometheus will only ingest the classic format even if the app exposes native histograms.

### 3. Update your queries

Native histograms simplify PromQL because you don't need `rate()` wrapping for quantiles:

| Classic | Native |
|---------|--------|
| `histogram_quantile(0.95, rate(..._bucket[5m]))` | `histogram_quantile(0.95, metric_name)` |
| `rate(..._sum[5m]) / rate(..._count[5m])` | `histogram_avg(metric_name)` |
| `rate(..._count[5m])` | `histogram_count(metric_name)` |
| N/A | `histogram_fraction(lower, upper, metric_name)` |

### 4. Migration strategy

You can run both formats simultaneously by keeping `Buckets` and adding the `NativeHistogram*` fields. The app will expose both, and Prometheus will ingest whichever format the scrape config allows. This lets you validate native histogram accuracy against your existing classic dashboards before removing the classic buckets.

[^1]: For an exponential distribution with mean μ, the CDF is `F(x) = 1 − e^(−x/μ)`. Solving `F(x) = 0.95`: `x = −μ × ln(0.05)`. For `/api/users` (μ = 0.1s): `x ≈ 0.2996s`.
[^2]: Native histograms use a exponential bucketing scheme where bucket boundaries grow by this factor. A factor of `1.1` produces boundaries at `..., 0.1, 0.11, 0.121, 0.1331, ...`. The factor determines the relative resolution error: a factor of `1.1` means any observation is at most ~10% away from a bucket boundary. Setting this to `1.0` disables native histograms entirely.
[^3]: When the number of populated buckets exceeds this limit, the histogram merges adjacent buckets (doubling the effective bucket factor) to reduce the count. This acts as a safety valve against high-cardinality distributions consuming too much memory. A value of `0` means no limit.
[^4]: After a bucket merge (caused by exceeding `MaxBucketNumber`), the histogram will not reset its observations until at least this duration has elapsed. This prevents a cascade of resets under bursty load. The counter only starts after the merge event, not from histogram creation.
[^5]: Observations with an absolute value at or below this threshold are counted in a special zero bucket rather than a regular exponential bucket. This avoids the problem of exponential buckets approaching zero requiring infinitely many buckets. For example, with a threshold of `0.001`, any observation in `[-0.001, 0.001]` goes into the zero bucket.
[^6]: Measured by timing 10 `curl` requests to the Prometheus query API and averaging (excluding the first cold request). Reproduce with: `for i in $(seq 1 10); do curl -s -o /dev/null -w "%{time_total}\n" --get 'http://localhost:9090/api/v1/query' --data-urlencode 'query=histogram_quantile(0.95, sum by (handler) (demo_request_duration_seconds))'; done` and the same with `query=histogram_quantile(0.95, sum by (handler, le) (rate(demo_classic_request_duration_seconds_bucket[5m])))`.
[^7]: Series counts and TSDB percentages from `curl http://localhost:9090/api/v1/status/tsdb`. The `seriesCountByMetricName` field shows 5,184 `_bucket` + 432 `_sum` + 432 `_count` = 6,048 classic series vs 432 native series. Samples per scrape from `scrape_samples_scraped{job="demo-app"}` (6,528 total) minus native series (432) gives 6,096 float samples; native contributes 432 histogram samples.
[^8]: From the `memoryInBytesByLabelName` field in `curl http://localhost:9090/api/v1/status/tsdb`. The `le` label (44,124 bytes) is used exclusively by classic histogram `_bucket` series to identify bucket boundaries. Native histograms encode buckets within the sample itself and do not use this label.