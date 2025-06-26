package main

import (
	"fmt"
	"os"

	"btree"
)

func main() {
	fmt.Println("=== Simple B-Tree Demo ===")
	fmt.Println("This demo shows basic B-Tree operations: Insert, Search, Delete")
	fmt.Println()

	// Clean up any existing files
	os.Remove("demo.gob")
	os.Remove("engine_demo.gob")
	defer func() {
		os.Remove("demo.gob")
		os.Remove("engine_demo.gob")
	}()

	// Part 1: Basic B-Tree Operations
	basicBTreeDemo()

	// Part 2: Persistent Operations with Engine
	engineDemo()

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("Key concepts demonstrated:")
	fmt.Println("✓ B-tree maintains sorted order automatically")
	fmt.Println("✓ Tree structure balances itself during insertions")
	fmt.Println("✓ Search operations are efficient (O(log n))")
	fmt.Println("✓ Persistence allows data to survive program restarts")
	fmt.Println("✓ Engine wrapper provides simple Put/Get/Delete interface")
}

func basicBTreeDemo() {
	fmt.Println("🌳 BASIC B-TREE OPERATIONS")
	fmt.Println("========================================")

	// Create a B-tree with order 3 (minimum degree = 2)
	bt := btree.New(3)
	fmt.Println("Created B-tree with order 3")

	// Insert some data
	fmt.Println("\n1. Inserting data:")
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
		bt.Insert(key, value)
		fmt.Printf("   Inserted: %s -> %s\n", key, value)
	}

	// Search for values
	fmt.Println("\n2. Searching for values:")
	searchKeys := []string{"apple", "fig", "nonexistent", "banana", "zzz"}
	for _, key := range searchKeys {
		if value, found := bt.Search(key); found {
			fmt.Printf("   %s -> %s ✓\n", key, value)
		} else {
			fmt.Printf("   %s -> NOT FOUND ✗\n", key)
		}
	}

	// Delete some values
	fmt.Println("\n3. Deleting values:")
	deleteKeys := []string{"cherry", "nonexistent", "grape"}
	for _, key := range deleteKeys {
		fmt.Printf("   Deleting: %s\n", key)
		bt.Delete(key)

		// Verify deletion
		if _, found := bt.Search(key); !found {
			fmt.Printf("     ✓ Successfully deleted\n")
		} else {
			fmt.Printf("     ✗ Still exists\n")
		}
	}

	// Search again after deletions
	fmt.Println("\n4. Searching after deletions:")
	for _, key := range []string{"cherry", "apple", "grape", "banana"} {
		if value, found := bt.Search(key); found {
			fmt.Printf("   %s -> %s ✓\n", key, value)
		} else {
			fmt.Printf("   %s -> NOT FOUND ✗\n", key)
		}
	}

	fmt.Println()
}

func engineDemo() {
	fmt.Println("💾 ENGINE (PERSISTENT) OPERATIONS")
	fmt.Println("========================================")

	// Create persistent engine
	engine, err := btree.Open("engine_demo.gob", 4)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created persistent B-tree engine with order 4")

	// Insert data using engine
	fmt.Println("\n1. Inserting data with persistence:")
	engineData := map[string]string{
		"user:1001":    "Alice Johnson",
		"user:1002":    "Bob Smith",
		"user:1003":    "Carol Davis",
		"order:O001":   "Order for user:1001",
		"order:O002":   "Order for user:1002",
		"product:P001": "Laptop Computer",
		"product:P002": "Wireless Mouse",
		"session:S001": "Active session for user:1001",
	}

	for key, value := range engineData {
		if err := engine.Put(key, value); err != nil {
			panic(err)
		}
		fmt.Printf("   Stored: %s -> %s\n", key, value)
	}

	// Read data
	fmt.Println("\n2. Reading data:")
	readKeys := []string{"user:1001", "order:O001", "product:P999", "session:S001"}
	for _, key := range readKeys {
		if value, found, err := engine.Get(key); err != nil {
			panic(err)
		} else if found {
			fmt.Printf("   %s -> %s ✓\n", key, value)
		} else {
			fmt.Printf("   %s -> NOT FOUND ✗\n", key)
		}
	}

	// Delete data
	fmt.Println("\n3. Deleting data:")
	if err := engine.Delete("order:O002"); err != nil {
		panic(err)
	}
	fmt.Println("   Deleted: order:O002")

	// Verify deletion
	if _, found, _ := engine.Get("order:O002"); !found {
		fmt.Println("   ✓ Deletion confirmed")
	}

	fmt.Println()
}
