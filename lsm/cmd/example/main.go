package main

import (
	"fmt"
	"lsm/bloom"
	"lsm/compaction"
	"lsm/lsmtree"
	"os"
	"path/filepath"
	"time"
)

func main() {
	fmt.Println("=== LSM Tree Comprehensive Learning Example ===")
	fmt.Println("This example demonstrates all aspects of LSM Tree implementation")
	fmt.Println()

	// Clean up any existing data
	dataDir := "lsm_data"
	os.RemoveAll(dataDir)

	// Part 1: Basic LSM Tree Operations
	fmt.Println("🌳 PART 1: Basic LSM Tree Operations")
	fmt.Println("=====================================")
	basicLSMDemo(dataDir)

	// Part 2: Enhanced Bloom Filter Features
	fmt.Println("\n🌸 PART 2: Bloom Filter Deep Dive")
	fmt.Println("==================================")
	bloomFilterDemo()

	// Part 3: Enhanced LSM with Compaction Strategies
	fmt.Println("\n🔧 PART 3: Advanced LSM with Compaction Strategies")
	fmt.Println("===================================================")
	enhancedLSMDemo(dataDir)

	// Part 4: Performance Analysis
	fmt.Println("\n📊 PART 4: Performance Analysis")
	fmt.Println("================================")
	performanceDemo(dataDir)

	fmt.Println("\n=== All LSM Tree Concepts Demonstrated ===")
	fmt.Println("✓ Basic operations: Put, Get, Compact")
	fmt.Println("✓ Memtable management and flushing")
	fmt.Println("✓ SSTable creation and management")
	fmt.Println("✓ Basic and enhanced Bloom filters")
	fmt.Println("✓ Multiple compaction strategies")
	fmt.Println("✓ Performance characteristics")
	fmt.Println("✓ Statistics and monitoring")
}

func basicLSMDemo(dataDir string) {
	// Create LSM tree with small threshold to see flushes in action
	fmt.Println("1. Creating basic LSM Tree...")
	tree, err := lsmtree.New(dataDir, 3) // Small threshold for demonstration
	if err != nil {
		panic(err)
	}
	fmt.Printf("   ✓ Created LSM tree with memtable threshold: 3\n")
	fmt.Printf("   ✓ Data directory: %s\n", dataDir)

	// Insert data to demonstrate memtable and SSTable creation
	fmt.Println("\n2. Inserting data (will trigger memtable flushes)...")
	data := map[string]string{
		"user:1001":    "Alice Johnson",
		"user:1002":    "Bob Smith",
		"user:1003":    "Carol Davis",
		"order:O001":   "Laptop",
		"order:O002":   "Mouse",
		"order:O003":   "Keyboard",
		"product:P001": "Electronics",
		"product:P002": "Accessories",
	}

	for key, value := range data {
		fmt.Printf("   Inserting: %s -> %s\n", key, value)
		if err := tree.Put(key, value); err != nil {
			panic(err)
		}
	}

	// Show created SSTable files
	fmt.Println("\n3. Checking created SSTable files...")
	listSSTables(dataDir)

	// Demonstrate reads
	fmt.Println("\n4. Reading data (searches memtable first, then SSTables)...")
	testKeys := []string{"user:1001", "order:O002", "product:P001", "nonexistent:key"}
	for _, key := range testKeys {
		if value, found, err := tree.Get(key); err != nil {
			panic(err)
		} else if found {
			fmt.Printf("   ✓ Found: %s -> %s\n", key, value)
		} else {
			fmt.Printf("   ✗ Not found: %s\n", key)
		}
	}

	// Demonstrate compaction
	fmt.Println("\n5. Running compaction (merges multiple SSTables into one)...")
	fmt.Printf("   Before compaction: ")
	countSSTables(dataDir)

	if err := tree.Compact(); err != nil {
		panic(err)
	}

	fmt.Printf("   After compaction: ")
	countSSTables(dataDir)

	// Verify data after compaction
	fmt.Println("\n6. Verifying data integrity after compaction...")
	for _, key := range []string{"user:1001", "order:O002", "product:P001"} {
		if value, found, err := tree.Get(key); err != nil {
			panic(err)
		} else if found {
			fmt.Printf("   ✓ Verified: %s -> %s\n", key, value)
		}
	}
}

func bloomFilterDemo() {
	fmt.Println("1. Basic Bloom Filter...")
	basicBloom := bloom.New(1000, 3)

	// Add some elements
	elements := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, elem := range elements {
		basicBloom.Add(elem)
		fmt.Printf("   Added: %s\n", elem)
	}

	// Test lookups
	fmt.Println("\n2. Testing basic Bloom filter lookups...")
	testElements := []string{"apple", "fig", "banana", "grape"}
	for _, elem := range testElements {
		if basicBloom.Contains(elem) {
			fmt.Printf("   %s: MAYBE in set\n", elem)
		} else {
			fmt.Printf("   %s: DEFINITELY NOT in set\n", elem)
		}
	}

	fmt.Println("\n3. Enhanced Bloom Filter with optimal sizing...")
	enhancedBloom := bloom.NewEnhanced(100, 0.01) // 1% false positive rate

	// Add elements and show statistics
	for _, elem := range elements {
		enhancedBloom.Add(elem)
	}

	stats := enhancedBloom.Stats()
	fmt.Printf("   Filter size: %d bits\n", stats.Size)
	fmt.Printf("   Hash functions: %d\n", stats.HashFunctions)
	fmt.Printf("   Fill ratio: %.2f%%\n", stats.FillRatio*100)
	fmt.Printf("   Expected FP rate: %.3f%%\n", stats.FalsePositiveRate*100)

	fmt.Println("\n4. Counting Bloom Filter (supports deletion)...")
	countingBloom := bloom.NewCounting(1000, 3)

	// Add elements
	for _, elem := range elements {
		countingBloom.Add(elem)
		fmt.Printf("   Added: %s\n", elem)
	}

	// Test deletion
	fmt.Printf("\n   Before deletion - contains 'banana': %t\n", countingBloom.Contains("banana"))
	countingBloom.Remove("banana")
	fmt.Printf("   After deletion - contains 'banana': %t\n", countingBloom.Contains("banana"))
}

func enhancedLSMDemo(dataDir string) {
	fmt.Println("1. Creating Enhanced LSM Tree with different compaction strategies...")

	strategies := []compaction.Strategy{
		compaction.NewSizeTieredStrategy(),
		compaction.NewLeveledStrategy(),
		compaction.NewTimeBasedStrategy(),
	}

	for i, strategy := range strategies {
		fmt.Printf("\n   Testing %s Strategy:\n", strategy.Name())
		testDir := fmt.Sprintf("%s_strategy_%d", dataDir, i)
		os.RemoveAll(testDir)

		// Create LSM tree with strategy
		tree, err := lsmtree.NewWithStrategy(testDir, 3, strategy)
		if err != nil {
			panic(err)
		}

		// Insert data in batches to trigger compaction
		batches := []map[string]string{
			{"batch1_key1": "value1", "batch1_key2": "value2"},
			{"batch2_key1": "value1", "batch2_key2": "value2"},
			{"batch3_key1": "value1", "batch3_key2": "value2"},
		}

		for j, batch := range batches {
			fmt.Printf("     Batch %d: ", j+1)
			for k, v := range batch {
				tree.Put(k, v)
				fmt.Printf("%s ", k)
			}
			fmt.Println()
		}

		// Show compaction info
		info := tree.GetCompactionInfo()
		fmt.Printf("     Should compact: %t\n", info.ShouldCompact)
		fmt.Printf("     Tables: %d, Selected: %d\n", info.TableCount, info.SelectedCount)

		// Show statistics
		if stats := tree.Stats(); stats != nil {
			fmt.Printf("     Writes: %d, Reads: %d\n", stats.TotalWrites, stats.TotalReads)
			fmt.Printf("     Compactions: %d, Flushes: %d\n", stats.CompactionCount, stats.TotalFlushes)
		} else {
			fmt.Printf("     Statistics not enabled for this tree\n")
		}
	}
}

func performanceDemo(dataDir string) {
	fmt.Println("1. Performance characteristics comparison...")

	// Test with different memtable sizes
	thresholds := []int{5, 10, 20}

	for _, threshold := range thresholds {
		fmt.Printf("\n   Testing with memtable threshold: %d\n", threshold)
		testDir := fmt.Sprintf("%s_perf_%d", dataDir, threshold)
		os.RemoveAll(testDir)

		tree, err := lsmtree.New(testDir, threshold)
		if err != nil {
			panic(err)
		}

		// Measure write performance
		start := time.Now()
		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("perf_key_%03d", i)
			value := fmt.Sprintf("perf_value_%03d", i)
			tree.Put(key, value)
		}
		writeTime := time.Since(start)

		// Measure read performance
		start = time.Now()
		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("perf_key_%03d", i)
			tree.Get(key)
		}
		readTime := time.Since(start)

		fmt.Printf("     Write time: %v (%.2f μs/op)\n", writeTime, float64(writeTime.Nanoseconds())/50/1000)
		fmt.Printf("     Read time: %v (%.2f μs/op)\n", readTime, float64(readTime.Nanoseconds())/50/1000)

		// Show file count
		fmt.Printf("     SSTable files: ")
		countSSTables(testDir)
	}
}

func listSSTables(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.sst"))
	if err != nil {
		fmt.Printf("   Error listing files: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("   No SSTable files found")
		return
	}

	fmt.Printf("   Found %d SSTable file(s):\n", len(files))
	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			fmt.Printf("     - %s (%d bytes)\n", filepath.Base(file), info.Size())
		}
	}
}

func countSSTables(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*.sst"))
	if err != nil {
		fmt.Printf("Error counting files: %v\n", err)
		return
	}
	fmt.Printf("%d SSTable file(s)\n", len(files))
}
