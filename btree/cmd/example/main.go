package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"btree"
)

func main() {
	fmt.Println("=== B-Tree Comprehensive Learning Example ===")
	fmt.Println("This example demonstrates all aspects of B-Tree implementation")
	fmt.Println()

	// Clean up any existing files
	cleanupFiles()

	// Part 1: Basic B-Tree Operations
	fmt.Println("🌳 PART 1: Basic B-Tree Operations")
	fmt.Println("===================================")
	basicBTreeDemo()

	// Part 2: Persistent Operations with Engine
	fmt.Println("\n💾 PART 2: Persistent Operations with Engine")
	fmt.Println("=============================================")
	persistentEngineDemo()

	// Part 3: Tree Structure Visualization
	fmt.Println("\n🔍 PART 3: Tree Structure & Growth Patterns")
	fmt.Println("============================================")
	structureVisualizationDemo()

	// Part 4: Performance Analysis
	fmt.Println("\n📊 PART 4: Performance Analysis")
	fmt.Println("================================")
	performanceAnalysisDemo()

	// Part 5: Advanced Features & Comparisons
	fmt.Println("\n🚀 PART 5: Advanced Features & Comparisons")
	fmt.Println("===========================================")
	advancedFeaturesDemo()

	fmt.Println("\n=== All B-Tree Concepts Demonstrated ===")
	fmt.Println("✓ Basic operations: Insert, Search, Delete")
	fmt.Println("✓ Balanced tree structure and node splitting/merging")
	fmt.Println("✓ Different tree orders and their effects")
	fmt.Println("✓ Persistence and recovery mechanisms")
	fmt.Println("✓ Performance characteristics and optimization")
	fmt.Println("✓ Memory vs disk usage analysis")
	fmt.Println("✓ Workload pattern analysis")
	fmt.Println("✓ Comparison with other data structures")

	// Clean up
	cleanupFiles()
}

func basicBTreeDemo() {
	fmt.Println("1. Creating B-tree with different orders...")

	orders := []int{3, 4, 5}
	for _, order := range orders {
		fmt.Printf("\n   Order %d (max %d keys per node):\n", order, 2*order-1)
		bt := btree.New(order)

		// Insert sample data
		data := []string{"M", "F", "P", "C", "A", "D", "Z", "E"}
		for _, key := range data {
			bt.Insert(key, fmt.Sprintf("value_%s", key))
		}

		// Test searches
		fmt.Printf("     Inserted: %s\n", strings.Join(data, ", "))
		testKeys := []string{"A", "M", "Z", "X"}
		for _, key := range testKeys {
			if _, found := bt.Search(key); found {
				fmt.Printf("     ✓ Found: %s\n", key)
			} else {
				fmt.Printf("     ✗ Not found: %s\n", key)
			}
		}

		// Test deletion
		bt.Delete("F")
		if _, found := bt.Search("F"); !found {
			fmt.Printf("     ✓ Successfully deleted: F\n")
		}
	}
}

func persistentEngineDemo() {
	fmt.Println("1. Creating persistent B-tree engine...")

	engine, err := btree.Open("comprehensive_demo.gob", 4)
	if err != nil {
		panic(err)
	}
	fmt.Printf("   ✓ Created persistent engine with order 4\n")

	// Insert structured data
	fmt.Println("\n2. Inserting structured data...")
	data := map[string]string{
		"user:1001":      "Alice Johnson",
		"user:1002":      "Bob Smith",
		"user:1003":      "Carol Davis",
		"order:O001":     "Laptop Order - $1200",
		"order:O002":     "Mouse Order - $25",
		"product:P001":   "Laptop Computer",
		"product:P002":   "Wireless Mouse",
		"config:timeout": "30s",
		"config:retries": "3",
	}

	for key, value := range data {
		if err := engine.Put(key, value); err != nil {
			panic(err)
		}
		fmt.Printf("   Stored: %s -> %s\n", key, value)
	}

	// Demonstrate range-like queries
	fmt.Println("\n3. Demonstrating prefix searches...")
	prefixes := []string{"user:", "order:", "config:"}
	for _, prefix := range prefixes {
		fmt.Printf("   Keys starting with '%s':\n", prefix)
		count := 0
		for key := range data {
			if strings.HasPrefix(key, prefix) {
				if value, found, _ := engine.Get(key); found {
					fmt.Printf("     %s -> %s\n", key, value)
					count++
				}
			}
		}
		fmt.Printf("     Found %d entries\n", count)
	}

	// Test persistence
	fmt.Println("\n4. Testing persistence across restarts...")

	// Simulate restart by creating new engine instance
	engine2, err := btree.Open("comprehensive_demo.gob", 4)
	if err != nil {
		panic(err)
	}

	// Verify data survived
	testKey := "user:1001"
	if value, found, _ := engine2.Get(testKey); found {
		fmt.Printf("   ✓ Data survived restart: %s -> %s\n", testKey, value)
	} else {
		fmt.Printf("   ✗ Data lost after restart\n")
	}
}

func structureVisualizationDemo() {
	fmt.Println("1. Demonstrating tree growth patterns...")

	bt := btree.New(3) // Order 3 for clear visualization
	fmt.Printf("   B-tree with order 3 (max 4 keys per node)\n")

	// Insert data step by step
	insertions := []string{"M", "F", "P", "C", "A", "D", "Z", "E", "K", "L"}

	for i, key := range insertions {
		bt.Insert(key, fmt.Sprintf("value_%s", key))
		if i%3 == 2 || i == len(insertions)-1 { // Show every 3rd insertion
			fmt.Printf("\n   After inserting %d keys (%s):\n", i+1, strings.Join(insertions[:i+1], ", "))
			visualizeTree(bt)
		}
	}

	fmt.Println("\n2. Comparing different tree orders...")
	orders := []int{3, 4, 5}
	numElements := 10

	for _, order := range orders {
		fmt.Printf("\n   Order %d with %d elements:\n", order, numElements)
		bt := btree.New(order)
		for i := 1; i <= numElements; i++ {
			key := fmt.Sprintf("K%02d", i)
			bt.Insert(key, fmt.Sprintf("V%02d", i))
		}
		visualizeTree(bt)
	}
}

func performanceAnalysisDemo() {
	fmt.Println("1. Performance comparison across different orders...")

	orders := []int{3, 4, 5, 10}
	dataSize := 1000

	fmt.Printf("%-8s %-12s %-12s %-12s %-12s\n", "Order", "Insert(μs)", "Search(μs)", "Delete(μs)", "Memory")
	fmt.Println("----------------------------------------------------------------")

	for _, order := range orders {
		bt := btree.New(order)

		// Generate test data
		keys := make([]string, dataSize)
		for i := 0; i < dataSize; i++ {
			keys[i] = fmt.Sprintf("key%06d", i)
		}

		// Measure insertion
		start := time.Now()
		for i, key := range keys {
			bt.Insert(key, fmt.Sprintf("value%06d", i))
		}
		insertTime := time.Since(start)

		// Measure search
		start = time.Now()
		for i := 0; i < 100; i++ {
			bt.Search(keys[rand.Intn(len(keys))])
		}
		searchTime := time.Since(start)

		// Measure deletion
		start = time.Now()
		for i := 0; i < 100; i++ {
			bt.Delete(keys[i])
		}
		deleteTime := time.Since(start)

		// Estimate memory usage
		memoryEstimate := estimateMemoryUsage(order, dataSize-100)

		fmt.Printf("%-8d %-12.2f %-12.2f %-12.2f %-12s\n",
			order,
			float64(insertTime.Nanoseconds())/float64(dataSize)/1000,
			float64(searchTime.Nanoseconds())/100/1000,
			float64(deleteTime.Nanoseconds())/100/1000,
			formatBytes(memoryEstimate))
	}

	fmt.Println("\n2. Workload pattern analysis...")
	testWorkloadPatterns()
}

func advancedFeaturesDemo() {
	fmt.Println("1. Memory vs Disk usage analysis...")

	dataSize := 500
	bt := btree.New(4)

	// Insert test data
	fmt.Printf("   Inserting %d key-value pairs...\n", dataSize)
	for i := 0; i < dataSize; i++ {
		key := fmt.Sprintf("user:%06d", i)
		value := fmt.Sprintf("User %d - %s", i, generateRandomString(30))
		bt.Insert(key, value)
	}

	// Save to disk and measure
	filename := "memory_analysis.gob"
	if err := bt.Save(filename); err != nil {
		panic(err)
	}

	if info, err := os.Stat(filename); err == nil {
		fmt.Printf("   Disk usage: %s\n", formatBytes(int(info.Size())))
		fmt.Printf("   Average per entry: %.2f bytes\n", float64(info.Size())/float64(dataSize))
	}

	memUsage := estimateMemoryUsage(4, dataSize)
	fmt.Printf("   Estimated memory usage: %s\n", formatBytes(memUsage))

	fmt.Println("\n2. B-Tree vs LSM Tree comparison...")
	showComparison()
}

// Helper functions
func visualizeTree(bt *btree.BTree) {
	// Since we can't access internal structure, show conceptual representation
	var existingKeys []string

	// Test common patterns to find existing keys
	testKeys := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"K01", "K02", "K03", "K04", "K05", "K06", "K07", "K08", "K09", "K10",
	}

	for _, key := range testKeys {
		if _, found := bt.Search(key); found {
			existingKeys = append(existingKeys, key)
		}
	}

	if len(existingKeys) == 0 {
		fmt.Println("     (empty tree)")
		return
	}

	fmt.Printf("     Keys in tree: %s\n", strings.Join(existingKeys, ", "))
	fmt.Printf("     Tree height: ~%.0f levels\n", estimateTreeHeight(len(existingKeys), 3))
}

func testWorkloadPatterns() {
	patterns := []struct {
		name        string
		description string
		testFunc    func() time.Duration
	}{
		{"Sequential Insert", "Insert keys in order", func() time.Duration {
			bt := btree.New(4)
			start := time.Now()
			for i := 0; i < 500; i++ {
				bt.Insert(fmt.Sprintf("key%06d", i), fmt.Sprintf("value%06d", i))
			}
			return time.Since(start)
		}},
		{"Random Insert", "Insert keys randomly", func() time.Duration {
			bt := btree.New(4)
			keys := make([]int, 500)
			for i := range keys {
				keys[i] = i
			}
			rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

			start := time.Now()
			for _, i := range keys {
				bt.Insert(fmt.Sprintf("key%06d", i), fmt.Sprintf("value%06d", i))
			}
			return time.Since(start)
		}},
		{"Range Search", "Search consecutive keys", func() time.Duration {
			bt := btree.New(4)
			for i := 0; i < 500; i++ {
				bt.Insert(fmt.Sprintf("key%06d", i), fmt.Sprintf("value%06d", i))
			}

			start := time.Now()
			for i := 100; i < 200; i++ {
				bt.Search(fmt.Sprintf("key%06d", i))
			}
			return time.Since(start)
		}},
		{"Random Search", "Search random keys", func() time.Duration {
			bt := btree.New(4)
			for i := 0; i < 500; i++ {
				bt.Insert(fmt.Sprintf("key%06d", i), fmt.Sprintf("value%06d", i))
			}

			start := time.Now()
			for i := 0; i < 100; i++ {
				bt.Search(fmt.Sprintf("key%06d", rand.Intn(500)))
			}
			return time.Since(start)
		}},
	}

	fmt.Printf("   %-20s %-25s %-15s\n", "Pattern", "Description", "Time")
	fmt.Println("   ------------------------------------------------------------")

	for _, pattern := range patterns {
		duration := pattern.testFunc()
		fmt.Printf("   %-20s %-25s %-15v\n", pattern.name, pattern.description, duration)
	}
}

func showComparison() {
	comparison := []struct {
		aspect string
		btree  string
		lsm    string
	}{
		{"Write Performance", "O(log n) - slower", "O(1) amortized - faster"},
		{"Read Performance", "O(log n) - consistent", "O(log n) - variable"},
		{"Space Efficiency", "Good - no duplicates", "Poor - duplicates until compaction"},
		{"Range Queries", "Excellent - sorted order", "Good - but may span multiple SSTables"},
		{"Memory Usage", "Moderate - tree structure", "Low - append-only writes"},
		{"Persistence", "Full tree serialization", "Incremental SSTable files"},
		{"Use Cases", "OLTP, consistent reads", "OLAP, write-heavy workloads"},
	}

	fmt.Printf("   %-20s %-25s %-30s\n", "Aspect", "B-Tree", "LSM Tree")
	fmt.Println("   -----------------------------------------------------------------------")

	for _, comp := range comparison {
		fmt.Printf("   %-20s %-25s %-30s\n", comp.aspect, comp.btree, comp.lsm)
	}

	fmt.Println("\n   Key Takeaways:")
	fmt.Println("   • B-Trees excel at read-heavy workloads with consistent performance")
	fmt.Println("   • LSM Trees excel at write-heavy workloads with high throughput")
	fmt.Println("   • B-Trees are better for OLTP (transactional) systems")
	fmt.Println("   • LSM Trees are better for OLAP (analytical) systems")
}

func estimateMemoryUsage(order, numKeys int) int {
	// Rough estimation: each node can hold 2*order-1 keys
	maxKeysPerNode := 2*order - 1
	estimatedNodes := numKeys / maxKeysPerNode
	if numKeys%maxKeysPerNode != 0 {
		estimatedNodes++
	}

	// Estimate bytes per node (keys + values + pointers)
	bytesPerNode := maxKeysPerNode*(20+30) + order*8 // rough estimate
	return estimatedNodes * bytesPerNode
}

func estimateTreeHeight(numKeys, order int) float64 {
	if numKeys <= 1 {
		return 1
	}
	maxKeysPerNode := 2*order - 1
	return 1 + estimateTreeHeight(numKeys/maxKeysPerNode, order)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}

func cleanupFiles() {
	files := []string{
		"comprehensive_demo.gob",
		"memory_analysis.gob",
		"demo.gob",
		"engine_demo.gob",
	}
	for _, file := range files {
		os.Remove(file)
	}
}
