# Infrastructure Learning Repository

This repository contains clean, educational implementations of fundamental data structures used in modern databases and storage systems. The implementations are designed for learning with simple examples and clear documentation.

## 📚 Learning Objectives

By exploring this repository, you will learn:

- **LSM Trees**: How write-optimized storage systems work (used in Cassandra, RocksDB, LevelDB)
- **B-Trees**: How traditional database indexes work (used in MySQL, PostgreSQL, SQLite)
- **Trade-offs**: When to use each data structure based on workload patterns
- **Implementation details**: Practical considerations for building storage systems

## 🌳 LSM Tree Implementation

The `lsm/` directory contains a unified LSM tree implementation with optional advanced features:

### Core Components
* **`memtable`** – In-memory buffer for recent writes (sorted map)
* **`bloom`** – Bloom filter implementations for efficient lookups
* **`sstable`** – Immutable sorted files on disk with Bloom filters
* **`lsmtree`** – Unified LSM tree with optional compaction strategies and statistics
* **`compaction`** – Multiple compaction strategies (Size-Tiered, Leveled, Time-Based)

### Features
- **Basic Operations**: Put, Get, Compact with clear explanations
- **Optional Advanced Features**: Statistics tracking and pluggable compaction strategies
- **Bloom Filters**: Efficient false-positive filtering for disk reads
- **Persistence**: Automatic loading of existing SSTables on startup

### Quick Start

```bash
# Simple demo - basic operations
cd lsm && go run ./cmd/demo

# Comprehensive example - all features including compaction strategies
cd lsm && go run ./cmd/example

# Run tests
cd lsm && go test ./lsmtree
```

**Key files created**: SSTables in data directories (automatically cleaned up in demos)

## 🌲 B-Tree Implementation

The `btree/` directory contains a clean B-tree implementation with persistence:

### Core Components
* **`btree.go`** – Core B-tree with insert, search, delete operations
* **`engine.go`** – Storage engine wrapper providing Put/Get/Delete interface
* **Persistence** – Automatic saving/loading using Go's gob encoding

### Features
- **Basic Operations**: Insert, Search, Delete with automatic balancing
- **Multiple Orders**: Support for different B-tree orders (branching factors)
- **Persistence**: Data survives program restarts via disk storage
- **Engine Interface**: Simple Put/Get/Delete API for easy integration

### Quick Start

```bash
# Simple demo - basic operations
cd btree && go run ./cmd/demo

# Comprehensive example - all features including persistence and performance
cd btree && go run ./cmd/example

# Run tests
cd btree && go test
```

**Key files created**: `.gob` files for persistence (automatically cleaned up in demos)

## ⚡ Performance Comparison & Benchmarks

The `benchmark/` directory contains comprehensive performance tests and integration tests:

### Running Benchmarks
```bash
# Performance benchmarks
cd benchmark && go test -bench . -benchmem

# Integration tests
cd benchmark && go test -v
```

### Expected Performance Characteristics

**LSM Trees excel at:**
- High write throughput (writes go to memory first)
- Write-heavy workloads
- Time-series data and logging systems

**B-Trees excel at:**
- Fast point lookups (balanced tree structure)
- Read-heavy workloads
- Range queries and sorted iteration

### Sample Benchmark Results
```
BenchmarkWriteLSM-3     121      8770492 ns/op      3105421 B/op    40414 allocs/op
BenchmarkWriteBTree-3     1  28951579208 ns/op  13945980560 B/op  50513455 allocs/op
BenchmarkReadLSM-3        1  1012496209 ns/op    561200192 B/op  12894572 allocs/op
BenchmarkReadBTree-3    668    1665341 ns/op             0 B/op        0 allocs/op
```

**Key takeaways:**
- LSM trees are ~3000x faster for writes (buffered in memory)
- B-trees are ~600x faster for reads (direct tree traversal)
- Choose based on your read/write ratio and consistency requirements

## 🧪 Testing Your Understanding

After running the comprehensive examples, try these exercises to deepen your understanding:

### LSM Tree Exercises
1. **Modify compaction strategies** in the enhanced LSM demo and observe their effects
2. **Experiment with Bloom filter parameters** and measure false positive rates
3. **Add more data** after compaction and see how the system handles mixed old/new data
4. **Implement a simple range query** that reads multiple consecutive keys
5. **Compare performance** with different memtable thresholds

### B-Tree Exercises
1. **Change the B-tree order** and observe how it affects tree height and performance
2. **Implement range scans** that return all keys between two values
3. **Add timing measurements** to compare search performance with different orders
4. **Test with different data patterns** (sequential vs random insertions)
5. **Analyze memory usage** patterns with different tree configurations

### Comparison Exercises
1. **Run the benchmarks** and analyze the performance trade-offs
2. **Design a hybrid system** that uses both structures for different data types
3. **Consider real-world scenarios**: When would you choose each structure?
4. **Implement a simple caching layer** on top of either structure

## 🔧 Quick Start Guide

```bash
# Clone and explore
git clone <repository-url>
cd infra-learning

# Try LSM Tree (simple demo)
cd lsm && go run ./cmd/demo

# Try B-Tree (simple demo)
cd ../btree && go run ./cmd/demo

# Run comprehensive examples
cd ../lsm && go run ./cmd/example
cd ../btree && go run ./cmd/example

# Run tests and benchmarks
cd ../lsm && go test ./lsmtree
cd ../btree && go test
cd ../benchmark && go test -v && go test -bench . -benchmem
```

## 📖 Advanced Topics (Optional)

The repository focuses on core concepts, but here are areas for further exploration:

- **Concurrent access**: Thread-safety and locking strategies
- **Write-ahead logging**: Durability and crash recovery
- **Compression**: Reducing storage overhead in SSTables
- **Leveled compaction**: More sophisticated LSM compaction strategies
- **B+ trees**: Leaf-linked B-trees for better range queries
- **Distributed systems**: Partitioning and replication strategies

## 📚 Additional Resources

- [The Log-Structured Merge-Tree (LSM-Tree)](http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.44.2782) - Original LSM paper
- [B-Trees and Database Indexes](https://use-the-index-luke.com/sql/anatomy/the-tree) - Comprehensive B-tree guide
- [Designing Data-Intensive Applications](https://dataintensive.net/) - Chapter 3 covers storage engines

## 🚀 Additional Services

The `services/` directory contains additional infrastructure examples:

- **`go-api/`**: HTTP API with OpenTelemetry instrumentation and Prometheus metrics
- **`otel-collector/`**: OpenTelemetry collector configuration
- **`docker-compose.yml`**: Complete observability stack with ClickHouse

These demonstrate how the data structures might be used in production systems with proper monitoring and observability.
