#!/bin/sh
# p95 latency per GET route, measured at the server (NFR-1: under 200 ms for lookup and
# phrase requests). Reads the API's Prometheus metrics from a URL or stdin:
#   scripts/p95.sh http://localhost:9090/metrics
#   docker compose exec api wget -qO- localhost:9090/metrics | scripts/p95.sh
# Interpolates within buckets as Prometheus's histogram_quantile does. Counts run from the
# API's start; with Prometheus, use histogram_quantile over a rate window instead.
set -eu
if [ $# -gt 0 ]; then curl --fail --silent "$1"; else cat; fi | awk '
/^vernacular_http_request_duration_seconds_bucket\{/ && /method="GET"/ {
  match($0, /route="[^"]*"/); route = substr($0, RSTART + 7, RLENGTH - 8)
  match($0, /le="[^"]*"/); le = substr($0, RSTART + 4, RLENGTH - 5)
  n[route, le] += $NF; les[le] = 1; routes[route] = 1
}
END {
  # httpx.latencyBuckets keeps an edge at 0.2 s, the NFR-1 target.
  if (!("0.2" in les)) { print "no le=\"0.2\" bucket: check latencyBuckets in httpx/metrics.go" > "/dev/stderr"; exit 1 }
  k = 0; for (le in les) if (le != "+Inf") edges[++k] = le + 0
  for (i = 2; i <= k; i++) for (j = i; j > 1 && edges[j - 1] > edges[j]; j--) { t = edges[j]; edges[j] = edges[j - 1]; edges[j - 1] = t }
  printf "%-34s %8s %9s %11s\n", "route", "requests", "p95 (ms)", "under 200ms"
  for (r in routes) {
    total = n[r, "+Inf"]; if (total == 0) continue
    want = 0.95 * total; p95 = -1; prev = 0; lo = 0
    for (i = 1; i <= k; i++) {
      c = n[r, edges[i] ""]
      if (p95 < 0 && c >= want) p95 = (c == prev) ? edges[i] : lo + (edges[i] - lo) * (want - prev) / (c - prev)
      prev = c; lo = edges[i]
    }
    p95txt = (p95 < 0) ? ("> " edges[k] * 1000) : sprintf("%.0f", p95 * 1000)
    printf "%-34s %8d %9s %10.1f%%\n", r, total, p95txt, 100 * n[r, "0.2"] / total
  }
}'
