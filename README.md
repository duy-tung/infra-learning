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
