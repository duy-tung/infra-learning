package main

import (
	"fmt"
	"os"
	"path/filepath"

	"lsm/lsmtree"
)

func main() {
	// Clean up any existing data
	dataDir := "demo_data"
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir) // Clean up after demo

	fmt.Println("=== Simple LSM Tree Demo ===")
	fmt.Println("This demo shows basic LSM Tree operations: Put, Get, and Compact")
	fmt.Println()

	// Create basic LSM tree with small threshold for demonstration
	tree, err := lsmtree.New(dataDir, 3) // Small threshold to trigger flushes
	if err != nil {
		panic(err)
	}
	fmt.Printf("✓ Created LSM tree with memtable threshold: 3\n")
	fmt.Printf("✓ Data directory: %s\n", dataDir)

	fmt.Println("\n1. Inserting data (will trigger memtable flushes)...")

	// Insert some data to trigger multiple flushes
	data := map[string]string{
		"apple":      "red fruit",
		"banana":     "yellow fruit",
		"cherry":     "red small fruit",
		"date":       "brown fruit",
		"elderberry": "purple berry",
		"fig":        "purple fruit",
		"grape":      "green/purple fruit",
		"honeydew":   "green melon",
	}

	for key, value := range data {
		fmt.Printf("   Inserting: %s -> %s\n", key, value)
		if err := tree.Put(key, value); err != nil {
			panic(err)
		}
	}

	fmt.Println("\n2. SSTable files created:")
	listSSTables(dataDir)

	fmt.Println("\n3. Reading values from LSM tree:")
	testKeys := []string{"apple", "fig", "nonexistent", "banana"}
	for _, key := range testKeys {
		if value, found, err := tree.Get(key); err != nil {
			panic(err)
		} else if found {
			fmt.Printf("   ✓ %s -> %s\n", key, value)
		} else {
			fmt.Printf("   ✗ %s -> NOT FOUND\n", key)
		}
	}

	fmt.Println("\n4. Before compaction:")
	listSSTables(dataDir)

	fmt.Println("\n5. Running compaction (merges all SSTables)...")
	if err := tree.Compact(); err != nil {
		panic(err)
	}

	fmt.Println("\n6. After compaction:")
	listSSTables(dataDir)

	fmt.Println("\n7. Verifying data integrity after compaction:")
	for _, key := range []string{"apple", "fig", "banana"} {
		if value, found, err := tree.Get(key); err != nil {
			panic(err)
		} else if found {
			fmt.Printf("   ✓ %s -> %s\n", key, value)
		}
	}

	fmt.Println("\n8. Adding more data after compaction...")
	newData := map[string]string{
		"kiwi":   "green fruit",
		"lemon":  "yellow citrus",
		"mango":  "orange tropical fruit",
		"orange": "orange citrus",
	}

	for key, value := range newData {
		fmt.Printf("   Inserting: %s -> %s\n", key, value)
		if err := tree.Put(key, value); err != nil {
			panic(err)
		}
	}

	fmt.Println("\n9. Final state:")
	listSSTables(dataDir)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Key concepts demonstrated:")
	fmt.Println("✓ Memtable fills up and flushes to SSTables")
	fmt.Println("✓ Multiple SSTables are created as data is inserted")
	fmt.Println("✓ Reads check memtable first, then SSTables (newest to oldest)")
	fmt.Println("✓ Compaction merges SSTables to reduce file count")
	fmt.Println("✓ Data integrity is maintained throughout operations")
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

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		fmt.Printf("   %s (%d bytes)\n", filepath.Base(file), info.Size())
	}
}
