# infra-learning

This repository contains learning experiments for infrastructure concepts.

## Go LSM Tree example

Under `lsm/` you will find a simple Log Structured Merge (LSM) tree
implementation written in Go. It is composed of the following packages:

* `memtable` – in-memory key/value store with flush support.
* `bloom` – very small Bloom filter implementation.
* `sstable` – read/write sorted string tables backed by text files.
* `lsmtree` – coordinates a memtable and a set of SSTables.
* `cmd/example` – example program exercising the tree.

### Build and run

From the repository root run:

```bash
go run ./lsm/cmd/example
```

Data files are written to `lsm/data` by default.

## B-tree implementation

The repository also contains a simple B-tree implementation under
`btree/`. It supports insertion, search and deletion with optional
persistence via gob files. A small storage engine wrapping the B-tree is
provided to offer `Put`, `Get` and `Delete` operations similar to the
LSM tree example.


## Benchmarks

A small benchmark under `benchmark/` compares read and write performance
between the LSM tree and B-tree using the same randomly generated data set.
Run it with:

```bash
go test -bench . -benchmem ./benchmark
```

Results on a sample run:

```
BenchmarkWriteLSM-3     121      8770492 ns/op      3105421 B/op    40414 allocs/op
BenchmarkWriteBTree-3     1  28951579208 ns/op  13945980560 B/op  50513455 allocs/op
BenchmarkReadLSM-3        1  1012496209 ns/op    561200192 B/op  12894572 allocs/op
BenchmarkReadBTree-3    668    1665341 ns/op             0 B/op        0 allocs/op
```

## HTTP API service

Under `services/go-api` there is a small Gin-based HTTP API. It is
instrumented using OpenTelemetry and exposes Prometheus metrics. When running the service you can
configure tracing via environment variables:

* `CLICKHOUSE_ENDPOINT` – ClickHouse connection endpoint.
* `CLICKHOUSE_DATABASE` – database name.
* `CLICKHOUSE_USERNAME` – authentication username.
* `CLICKHOUSE_PASSWORD` – authentication password.
* `OTEL_SERVICE_NAME` – service name (defaults to `go-api`).
* `OTEL_SERVICE_VERSION` – optional version tag.
* `OTEL_SAMPLER_RATIO` – sampling ratio between `0` and `1`.
* `OTEL_CUSTOM_TAGS` – comma separated custom span tags (`key=value`).

Example:

```bash
CLICKHOUSE_ENDPOINT=http://clickhouse:8123 \
CLICKHOUSE_DATABASE=otel \
CLICKHOUSE_USERNAME=default \
CLICKHOUSE_PASSWORD=secret \
OTEL_SERVICE_NAME=go-api \
OTEL_SERVICE_VERSION=1.0.0 \
OTEL_SAMPLER_RATIO=1 \
OTEL_CUSTOM_TAGS=env=dev,team=infra \
go run ./services/go-api
```

The service exposes a `/metrics` endpoint compatible with Prometheus. Histogram
buckets are configured so you can query p50, p90, p95 and p99 latencies.
