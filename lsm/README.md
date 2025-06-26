# LSM Tree Implementation

A clean, educational implementation of Log-Structured Merge (LSM) trees with optional advanced features.

## Overview

LSM trees are write-optimized data structures used in systems like Cassandra, RocksDB, and LevelDB. They excel at handling high write throughput by buffering writes in memory and periodically flushing to disk.

## Architecture

```
┌─────────────┐    ┌──────────────┐    ┌──────────────┐
│  Memtable   │───▶│   SSTable    │───▶│   SSTable    │
│ (In Memory) │    │ (On Disk #1) │    │ (On Disk #2) │
└─────────────┘    └──────────────┘    └──────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
   Put/Get              Get (with              Get
   Operations           Bloom Filter)      (Fallback)
```

## Core Components

- **Memtable**: In-memory sorted map that buffers recent writes
- **SSTable**: Immutable sorted files on disk with Bloom filters
- **Bloom Filters**: Probabilistic data structure to avoid unnecessary disk reads
- **Compaction**: Process to merge SSTables and reclaim space

## Usage

### Basic Usage

```go
// Create a basic LSM tree
tree, err := lsmtree.New("data_dir", 100) // threshold = 100 entries
if err != nil {
    panic(err)
}

// Insert data
tree.Put("key1", "value1")
tree.Put("key2", "value2")

// Read data
value, found, err := tree.Get("key1")
if found {
    fmt.Printf("Found: %s\n", value)
}

// Compact SSTables
tree.Compact()
```

### Advanced Usage with Compaction Strategies

```go
// Create LSM tree with compaction strategy and statistics
strategy := compaction.NewSizeTieredStrategy()
tree, err := lsmtree.NewWithStrategy("data_dir", 100, strategy)
if err != nil {
    panic(err)
}

// Operations are the same
tree.Put("key", "value")

// Access statistics
if stats := tree.Stats(); stats != nil {
    fmt.Printf("Total writes: %d\n", stats.TotalWrites)
    fmt.Printf("Bloom filter saves: %d\n", stats.BloomFilterSaves)
}

// Get compaction information
info := tree.GetCompactionInfo()
fmt.Printf("Strategy: %s, Should compact: %t\n", 
    info.Strategy, info.ShouldCompact)
```

## Available Compaction Strategies

1. **Size-Tiered**: Groups SSTables by similar sizes, good for write-heavy workloads
2. **Leveled**: Organizes SSTables in levels, better for read-heavy workloads  
3. **Time-Based**: Compacts based on SSTable age

## Running Examples

### Simple Demo
```bash
go run ./cmd/demo
```
Shows basic Put, Get, and Compact operations with clear explanations.

### Comprehensive Example
```bash
go run ./cmd/example
```
Demonstrates all features including:
- Basic operations
- Bloom filter functionality
- Different compaction strategies
- Performance analysis

## Running Tests

```bash
go test ./lsmtree
```

Tests cover:
- Basic operations (Put/Get/Compact)
- Compaction strategies and statistics
- Data persistence across restarts
- Overwrite behavior and data integrity

## Key Concepts

### Write Path
1. Data goes to memtable (in memory)
2. When memtable is full, flush to SSTable on disk
3. Reads check memtable first, then SSTables (newest to oldest)

### Read Path
1. Check memtable first (fastest)
2. For each SSTable (newest to oldest):
   - Check Bloom filter (avoid disk read if key definitely not present)
   - If Bloom filter says "maybe", read from disk

### Compaction
- Merges multiple SSTables into fewer, larger ones
- Removes duplicate/overwritten keys
- Different strategies optimize for different workloads

## Performance Characteristics

**Strengths:**
- Excellent write performance (O(1) for memtable writes)
- Good for write-heavy workloads
- Efficient space utilization after compaction

**Trade-offs:**
- Read performance depends on number of SSTables
- Compaction can cause temporary I/O spikes
- More complex than simpler data structures

## Educational Value

This implementation demonstrates:
- How modern NoSQL databases handle writes efficiently
- Trade-offs between write and read performance
- The role of Bloom filters in reducing disk I/O
- Different compaction strategies and their use cases
- Statistics collection for performance monitoring
