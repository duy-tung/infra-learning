package lsmtree

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lsm/compaction"
	"lsm/memtable"
	"lsm/sstable"
)

// LSMTree coordinates memtable and SSTables with optional advanced features.
type LSMTree struct {
	Mem    *memtable.Memtable
	Tables []*sstable.SSTable
	Dir    string
	nextID int

	// Optional advanced features
	strategy *compaction.Strategy // nil for basic mode
	stats    *LSMStats            // nil for basic mode
}

// LSMStats tracks performance metrics
type LSMStats struct {
	TotalWrites      uint64
	TotalReads       uint64
	MemtableHits     uint64
	SSTableHits      uint64
	BloomFilterSaves uint64
	CompactionCount  uint64
	TotalFlushes     uint64
}

// New creates a basic LSM tree without advanced features.
func New(dir string, threshold int) (*LSMTree, error) {
	return newLSMTree(dir, threshold, nil, false)
}

// NewWithStrategy creates an LSM tree with a compaction strategy and statistics tracking.
func NewWithStrategy(dir string, threshold int, strategy compaction.Strategy) (*LSMTree, error) {
	return newLSMTree(dir, threshold, &strategy, true)
}

// newLSMTree is the internal constructor that handles both basic and advanced modes.
func newLSMTree(dir string, threshold int, strategy *compaction.Strategy, enableStats bool) (*LSMTree, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	t := &LSMTree{
		Mem:      memtable.New(threshold),
		Dir:      dir,
		strategy: strategy,
	}

	if enableStats {
		t.stats = &LSMStats{}
	}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Ext(path) == ".sst" {
			table, err := sstable.Load(path)
			if err != nil {
				return err
			}
			t.Tables = append(t.Tables, table)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(t.Tables, func(i, j int) bool { return t.Tables[i].Path < t.Tables[j].Path })
	t.nextID = len(t.Tables)
	return t, nil
}

// Put inserts a key-value pair.
func (t *LSMTree) Put(key, value string) error {
	// Track statistics if enabled
	if t.stats != nil {
		t.stats.TotalWrites++
	}

	t.Mem.Put(key, value)
	if t.Mem.IsFull() {
		if t.stats != nil {
			t.stats.TotalFlushes++
		}

		if err := t.flush(); err != nil {
			return err
		}

		// Check if compaction is needed using strategy (if available)
		if t.strategy != nil && (*t.strategy).ShouldCompact(t.Tables) {
			return t.CompactWithStrategy()
		}
	}
	return nil
}

func (t *LSMTree) flush() error {
	kvs := t.Mem.Flush()
	path := filepath.Join(t.Dir, fmt.Sprintf("ss-%d.sst", t.nextID))
	tbl, err := sstable.New(path, kvs)
	if err != nil {
		return err
	}
	t.Tables = append(t.Tables, tbl)
	t.nextID++
	return nil
}

// Get searches memtable then SSTables newest to oldest.
func (t *LSMTree) Get(key string) (string, bool, error) {
	// Track statistics if enabled
	if t.stats != nil {
		t.stats.TotalReads++
	}

	// Check memtable first
	if v, ok := t.Mem.Get(key); ok {
		if t.stats != nil {
			t.stats.MemtableHits++
		}
		return v, true, nil
	}

	// Check SSTables newest to oldest
	for i := len(t.Tables) - 1; i >= 0; i-- {
		// Use Bloom filter to avoid unnecessary disk reads (if available)
		if t.Tables[i].Bloom != nil && !t.Tables[i].Bloom.Contains(key) {
			if t.stats != nil {
				t.stats.BloomFilterSaves++
			}
			continue
		}

		if v, ok, err := t.Tables[i].Get(key); err != nil {
			return "", false, err
		} else if ok {
			if t.stats != nil {
				t.stats.SSTableHits++
			}
			return v, true, nil
		}
	}

	return "", false, nil
}

// Compact merges all tables into one (basic compaction).
func (t *LSMTree) Compact() error {
	if len(t.Tables) < 2 {
		return nil
	}

	if t.stats != nil {
		t.stats.CompactionCount++
	}

	merged := make(map[string]string)
	// newer tables override old
	for _, tbl := range t.Tables {
		f, err := os.Open(tbl.Path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			parts := strings.SplitN(scanner.Text(), "\t", 2)
			if len(parts) == 2 {
				merged[parts[0]] = parts[1]
			}
		}
		f.Close()
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	// use memtable KV type to sort
	kvs := make([]memtable.KV, 0, len(merged))
	for k, v := range merged {
		kvs = append(kvs, memtable.KV{Key: k, Value: v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].Key < kvs[j].Key })
	path := filepath.Join(t.Dir, fmt.Sprintf("ss-%d.sst", t.nextID))
	tbl, err := sstable.New(path, kvs)
	if err != nil {
		return err
	}
	// remove old tables
	for _, old := range t.Tables {
		os.Remove(old.Path)
	}
	t.Tables = []*sstable.SSTable{tbl}
	t.nextID++
	return nil
}

// CompactWithStrategy uses the configured strategy for compaction.
func (t *LSMTree) CompactWithStrategy() error {
	if t.strategy == nil {
		// Fall back to basic compaction if no strategy is set
		return t.Compact()
	}

	if !(*t.strategy).ShouldCompact(t.Tables) {
		return nil
	}

	if t.stats != nil {
		t.stats.CompactionCount++
	}

	selectedTables := (*t.strategy).SelectTables(t.Tables)
	if len(selectedTables) < 2 {
		return nil // Nothing to compact
	}

	// Merge selected tables
	merged := make(map[string]string)

	// Process tables in order (newer tables override older ones)
	for _, tbl := range selectedTables {
		f, err := os.Open(tbl.Path)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			parts := strings.SplitN(scanner.Text(), "\t", 2)
			if len(parts) == 2 {
				merged[parts[0]] = parts[1]
			}
		}
		f.Close()

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	// Convert to sorted KV pairs
	kvs := make([]memtable.KV, 0, len(merged))
	for k, v := range merged {
		kvs = append(kvs, memtable.KV{Key: k, Value: v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].Key < kvs[j].Key })

	// Create new SSTable
	newPath := filepath.Join(t.Dir, fmt.Sprintf("ss-%d.sst", t.nextID))
	newTable, err := sstable.New(newPath, kvs)
	if err != nil {
		return err
	}

	// Remove old tables from list and disk
	var remainingTables []*sstable.SSTable
	selectedPaths := make(map[string]bool)
	for _, tbl := range selectedTables {
		selectedPaths[tbl.Path] = true
		os.Remove(tbl.Path)
	}

	for _, tbl := range t.Tables {
		if !selectedPaths[tbl.Path] {
			remainingTables = append(remainingTables, tbl)
		}
	}

	// Add new table and update state
	t.Tables = append(remainingTables, newTable)
	t.nextID++

	return nil
}

// Stats returns performance statistics (nil if statistics are not enabled).
func (t *LSMTree) Stats() *LSMStats {
	return t.stats
}

// SetStrategy changes the compaction strategy.
func (t *LSMTree) SetStrategy(strategy compaction.Strategy) {
	t.strategy = &strategy
}

// String returns a formatted string representation of stats.
func (s LSMStats) String() string {
	hitRate := float64(0)
	if s.TotalReads > 0 {
		hitRate = float64(s.MemtableHits+s.SSTableHits) / float64(s.TotalReads) * 100
	}

	bloomEfficiency := float64(0)
	if s.TotalReads > 0 {
		bloomEfficiency = float64(s.BloomFilterSaves) / float64(s.TotalReads) * 100
	}

	return fmt.Sprintf(`LSM Tree Statistics:
  Total Writes: %d
  Total Reads: %d
  Memtable Hits: %d
  SSTable Hits: %d
  Hit Rate: %.2f%%
  Bloom Filter Saves: %d (%.2f%% efficiency)
  Total Flushes: %d
  Compactions: %d`,
		s.TotalWrites, s.TotalReads, s.MemtableHits, s.SSTableHits,
		hitRate, s.BloomFilterSaves, bloomEfficiency,
		s.TotalFlushes, s.CompactionCount)
}

// CompactionInfo provides details about the current state.
type CompactionInfo struct {
	Strategy      string
	ShouldCompact bool
	TableCount    int
	TotalSize     int64
	SelectedCount int
}

// GetCompactionInfo returns information about compaction readiness.
func (t *LSMTree) GetCompactionInfo() *CompactionInfo {
	if t.strategy == nil {
		return &CompactionInfo{
			Strategy:      "Basic",
			ShouldCompact: len(t.Tables) >= 2,
			TableCount:    len(t.Tables),
			TotalSize:     t.getTotalSize(),
			SelectedCount: len(t.Tables),
		}
	}

	shouldCompact := (*t.strategy).ShouldCompact(t.Tables)
	selectedTables := (*t.strategy).SelectTables(t.Tables)

	return &CompactionInfo{
		Strategy:      (*t.strategy).Name(),
		ShouldCompact: shouldCompact,
		TableCount:    len(t.Tables),
		TotalSize:     t.getTotalSize(),
		SelectedCount: len(selectedTables),
	}
}

// getTotalSize calculates the total size of all SSTables.
func (t *LSMTree) getTotalSize() int64 {
	totalSize := int64(0)
	for _, table := range t.Tables {
		if info, err := os.Stat(table.Path); err == nil {
			totalSize += info.Size()
		}
	}
	return totalSize
}

// String returns formatted compaction information.
func (ci CompactionInfo) String() string {
	return fmt.Sprintf(`Compaction Info:
  Strategy: %s
  Should Compact: %t
  Total Tables: %d
  Total Size: %d bytes
  Selected for Compaction: %d tables`,
		ci.Strategy, ci.ShouldCompact, ci.TableCount, ci.TotalSize, ci.SelectedCount)
}
