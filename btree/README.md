# B-Tree Implementation

A clean, educational implementation of B-trees with persistence support.

## Overview

B-trees are balanced tree data structures used in databases and file systems like MySQL, PostgreSQL, and SQLite. They maintain sorted data and provide efficient search, insertion, and deletion operations.

## Architecture

```
                    [M | P]
                   /   |   \
              [C|F]   [K|L]  [S|X|Z]
             /  |  \   /  \   /  |  |  \
           [A] [D] [G] [J] [N] [Q] [U] [Y]
```

Each node can contain multiple keys and children, keeping the tree balanced and minimizing disk I/O.

## Core Components

- **BTree**: Core B-tree structure with configurable order (branching factor)
- **Node**: Individual tree nodes containing keys, values, and child pointers
- **Engine**: Persistent wrapper providing simple Put/Get/Delete interface

## Usage

### Basic B-Tree Operations

```go
// Create a B-tree with order 3 (max 2 keys per node)
bt := btree.New(3)

// Insert data
bt.Insert("apple", "red fruit")
bt.Insert("banana", "yellow fruit")
bt.Insert("cherry", "red small fruit")

// Search for data
value, found := bt.Search("apple")
if found {
    fmt.Printf("Found: %s\n", value)
}

// Delete data
bt.Delete("banana")

// Save to disk
err := bt.Save("my_btree.gob")
if err != nil {
    panic(err)
}

// Load from disk
loadedBT, err := btree.Load("my_btree.gob")
if err != nil {
    panic(err)
}
```

### Persistent Engine Interface

```go
// Create or open persistent B-tree
engine, err := btree.Open("database.gob", 4) // order 4
if err != nil {
    panic(err)
}

// Simple Put/Get/Delete interface
err = engine.Put("user:1001", "Alice Johnson")
if err != nil {
    panic(err)
}

value, found, err := engine.Get("user:1001")
if err != nil {
    panic(err)
}
if found {
    fmt.Printf("User: %s\n", value)
}

err = engine.Delete("user:1001")
if err != nil {
    panic(err)
}
```

## Tree Order (Branching Factor)

The order determines the maximum number of children each node can have:
- **Order 3**: Max 2 keys per node, max 3 children
- **Order 4**: Max 3 keys per node, max 4 children
- **Higher orders**: Fewer levels but larger nodes

### Choosing the Right Order

- **Lower orders (3-4)**: Good for learning, easier to visualize
- **Higher orders (100+)**: Better for real databases, fewer disk reads

## Running Examples

### Simple Demo
```bash
go run ./cmd/demo
```
Shows basic Insert, Search, Delete operations and persistence.

### Comprehensive Example
```bash
go run ./cmd/example
```
Demonstrates all features including:
- Basic operations with different orders
- Persistence and recovery
- Tree structure visualization
- Performance analysis

## Running Tests

```bash
go test
```

Tests cover:
- Basic CRUD operations
- Persistence functionality
- Tree balancing behavior
- Engine wrapper functionality

## Key Concepts

### Self-Balancing
- B-trees automatically maintain balance during insertions and deletions
- Nodes split when they become too full
- Nodes merge when they become too empty
- Tree height grows/shrinks only at the root

### Disk Efficiency
- Designed to minimize disk I/O operations
- Each node typically corresponds to a disk page
- Higher branching factor = fewer levels = fewer disk reads

### Sorted Order
- Keys are always maintained in sorted order
- Enables efficient range queries
- In-order traversal gives sorted sequence

## Performance Characteristics

**Time Complexity:**
- Search: O(log n)
- Insert: O(log n)  
- Delete: O(log n)

**Strengths:**
- Excellent read performance
- Maintains sorted order
- Self-balancing
- Efficient for range queries

**Trade-offs:**
- Writes require tree rebalancing
- More complex than hash tables
- Memory overhead for tree structure

## Comparison with LSM Trees

| Aspect | B-Tree | LSM Tree |
|--------|--------|----------|
| **Reads** | Fast (O(log n)) | Slower (check multiple files) |
| **Writes** | Slower (rebalancing) | Fast (append-only) |
| **Use Case** | Read-heavy workloads | Write-heavy workloads |
| **Examples** | MySQL, PostgreSQL | Cassandra, RocksDB |

## Educational Value

This implementation demonstrates:
- How traditional databases index data efficiently
- Self-balancing tree algorithms
- Trade-offs between read and write performance
- Persistence mechanisms in database systems
- The importance of choosing appropriate data structures
