package lsmtree

import (
	"os"
	"testing"

	"lsm/compaction"
)

func TestBasicLSMOperations(t *testing.T) {
	// Clean up test directory
	testDir := "test_basic_lsm"
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// Create basic LSM tree
	tree, err := New(testDir, 3)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}

	// Test Put operations
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4", // This should trigger a flush
	}

	for key, value := range testData {
		if err := tree.Put(key, value); err != nil {
			t.Fatalf("Failed to put %s: %v", key, err)
		}
	}

	// Test Get operations
	for key, expectedValue := range testData {
		value, found, err := tree.Get(key)
		if err != nil {
			t.Fatalf("Failed to get %s: %v", key, err)
		}
		if !found {
			t.Fatalf("Key %s not found", key)
		}
		if value != expectedValue {
			t.Fatalf("Expected %s, got %s for key %s", expectedValue, value, key)
		}
	}

	// Test non-existent key
	_, found, err := tree.Get("nonexistent")
	if err != nil {
		t.Fatalf("Failed to get nonexistent key: %v", err)
	}
	if found {
		t.Fatalf("Found nonexistent key")
	}

	// Test compaction
	if err := tree.Compact(); err != nil {
		t.Fatalf("Failed to compact: %v", err)
	}

	// Verify data after compaction
	for key, expectedValue := range testData {
		value, found, err := tree.Get(key)
		if err != nil {
			t.Fatalf("Failed to get %s after compaction: %v", key, err)
		}
		if !found {
			t.Fatalf("Key %s not found after compaction", key)
		}
		if value != expectedValue {
			t.Fatalf("Expected %s, got %s for key %s after compaction", expectedValue, value, key)
		}
	}
}

func TestLSMWithStrategy(t *testing.T) {
	// Clean up test directory
	testDir := "test_strategy_lsm"
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// Create LSM tree with strategy
	strategy := compaction.NewSizeTieredStrategy()
	tree, err := NewWithStrategy(testDir, 2, strategy)
	if err != nil {
		t.Fatalf("Failed to create LSM tree with strategy: %v", err)
	}

	// Verify statistics are enabled
	if tree.Stats() == nil {
		t.Fatalf("Statistics should be enabled for tree with strategy")
	}

	// Test Put operations
	testData := map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": "4",
		"e": "5",
	}

	for key, value := range testData {
		if err := tree.Put(key, value); err != nil {
			t.Fatalf("Failed to put %s: %v", key, err)
		}
	}

	// Check statistics
	stats := tree.Stats()
	if stats.TotalWrites != uint64(len(testData)) {
		t.Fatalf("Expected %d writes, got %d", len(testData), stats.TotalWrites)
	}

	// Test reads and verify statistics
	for key := range testData {
		_, found, err := tree.Get(key)
		if err != nil {
			t.Fatalf("Failed to get %s: %v", key, err)
		}
		if !found {
			t.Fatalf("Key %s not found", key)
		}
	}

	stats = tree.Stats()
	if stats.TotalReads != uint64(len(testData)) {
		t.Fatalf("Expected %d reads, got %d", len(testData), stats.TotalReads)
	}

	// Test compaction info
	info := tree.GetCompactionInfo()
	if info == nil {
		t.Fatalf("Compaction info should not be nil")
	}
	if info.Strategy != strategy.Name() {
		t.Fatalf("Expected strategy %s, got %s", strategy.Name(), info.Strategy)
	}
}

func TestLSMPersistence(t *testing.T) {
	// Clean up test directory
	testDir := "test_persistence_lsm"
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	testData := map[string]string{
		"persistent1": "value1",
		"persistent2": "value2",
		"persistent3": "value3",
	}

	// Create and populate LSM tree
	{
		tree, err := New(testDir, 2)
		if err != nil {
			t.Fatalf("Failed to create LSM tree: %v", err)
		}

		for key, value := range testData {
			if err := tree.Put(key, value); err != nil {
				t.Fatalf("Failed to put %s: %v", key, err)
			}
		}

		// Force flush any remaining data in memtable by adding one more item
		if err := tree.Put("flush_trigger", "dummy"); err != nil {
			t.Fatalf("Failed to put flush trigger: %v", err)
		}
	}

	// Create new LSM tree instance (simulating restart)
	{
		tree, err := New(testDir, 2)
		if err != nil {
			t.Fatalf("Failed to create LSM tree after restart: %v", err)
		}

		// Verify all data is still there
		for key, expectedValue := range testData {
			value, found, err := tree.Get(key)
			if err != nil {
				t.Fatalf("Failed to get %s after restart: %v", key, err)
			}
			if !found {
				t.Fatalf("Key %s not found after restart", key)
			}
			if value != expectedValue {
				t.Fatalf("Expected %s, got %s for key %s after restart", expectedValue, value, key)
			}
		}
	}
}

func TestLSMOverwrite(t *testing.T) {
	// Clean up test directory
	testDir := "test_overwrite_lsm"
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// Create LSM tree
	tree, err := New(testDir, 2)
	if err != nil {
		t.Fatalf("Failed to create LSM tree: %v", err)
	}

	// Insert initial value
	if err := tree.Put("key", "value1"); err != nil {
		t.Fatalf("Failed to put initial value: %v", err)
	}

	// Overwrite with new value
	if err := tree.Put("key", "value2"); err != nil {
		t.Fatalf("Failed to put overwrite value: %v", err)
	}

	// Verify we get the latest value
	value, found, err := tree.Get("key")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if !found {
		t.Fatalf("Key not found")
	}
	if value != "value2" {
		t.Fatalf("Expected value2, got %s", value)
	}

	// Compact and verify latest value is preserved
	if err := tree.Compact(); err != nil {
		t.Fatalf("Failed to compact: %v", err)
	}

	value, found, err = tree.Get("key")
	if err != nil {
		t.Fatalf("Failed to get key after compaction: %v", err)
	}
	if !found {
		t.Fatalf("Key not found after compaction")
	}
	if value != "value2" {
		t.Fatalf("Expected value2 after compaction, got %s", value)
	}
}
