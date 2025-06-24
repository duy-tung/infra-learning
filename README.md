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

