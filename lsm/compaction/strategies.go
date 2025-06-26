package compaction

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lsm/sstable"
)

// Strategy defines different compaction approaches
type Strategy interface {
	ShouldCompact(tables []*sstable.SSTable) bool
	SelectTables(tables []*sstable.SSTable) []*sstable.SSTable
	Name() string
}

// SizeTieredStrategy compacts tables of similar sizes
type SizeTieredStrategy struct {
	MinTables    int     // Minimum tables to trigger compaction
	SizeRatio    float64 // Size ratio threshold
	MaxTableSize int64   // Maximum size before forced compaction
}

func NewSizeTieredStrategy() *SizeTieredStrategy {
	return &SizeTieredStrategy{
		MinTables:    4,
		SizeRatio:    2.0,
		MaxTableSize: 1024 * 1024, // 1MB
	}
}

func (s *SizeTieredStrategy) Name() string {
	return "Size-Tiered"
}

func (s *SizeTieredStrategy) ShouldCompact(tables []*sstable.SSTable) bool {
	if len(tables) < s.MinTables {
		return false
	}

	// Group tables by size tiers
	tiers := s.groupBySize(tables)
	
	// Check if any tier has enough tables
	for _, tier := range tiers {
		if len(tier) >= s.MinTables {
			return true
		}
	}
	
	// Check for oversized tables
	for _, table := range tables {
		if s.getTableSize(table) > s.MaxTableSize {
			return true
		}
	}
	
	return false
}

func (s *SizeTieredStrategy) SelectTables(tables []*sstable.SSTable) []*sstable.SSTable {
	// Group by size and select the tier with most tables
	tiers := s.groupBySize(tables)
	
	var bestTier []*sstable.SSTable
	maxCount := 0
	
	for _, tier := range tiers {
		if len(tier) > maxCount {
			maxCount = len(tier)
			bestTier = tier
		}
	}
	
	return bestTier
}

func (s *SizeTieredStrategy) groupBySize(tables []*sstable.SSTable) [][]*sstable.SSTable {
	// Sort tables by size
	sorted := make([]*sstable.SSTable, len(tables))
	copy(sorted, tables)
	
	sort.Slice(sorted, func(i, j int) bool {
		return s.getTableSize(sorted[i]) < s.getTableSize(sorted[j])
	})
	
	var tiers [][]*sstable.SSTable
	var currentTier []*sstable.SSTable
	var lastSize int64
	
	for _, table := range sorted {
		size := s.getTableSize(table)
		
		if lastSize == 0 || float64(size)/float64(lastSize) <= s.SizeRatio {
			currentTier = append(currentTier, table)
		} else {
			if len(currentTier) > 0 {
				tiers = append(tiers, currentTier)
			}
			currentTier = []*sstable.SSTable{table}
		}
		lastSize = size
	}
	
	if len(currentTier) > 0 {
		tiers = append(tiers, currentTier)
	}
	
	return tiers
}

func (s *SizeTieredStrategy) getTableSize(table *sstable.SSTable) int64 {
	if info, err := os.Stat(table.Path); err == nil {
		return info.Size()
	}
	return 0
}

// LeveledStrategy implements leveled compaction
type LeveledStrategy struct {
	MaxLevel     int
	LevelSizes   []int64 // Max size for each level
	LevelRatios  []int   // Size multiplier between levels
}

func NewLeveledStrategy() *LeveledStrategy {
	return &LeveledStrategy{
		MaxLevel:    7,
		LevelSizes:  []int64{10 * 1024, 100 * 1024, 1024 * 1024}, // 10KB, 100KB, 1MB
		LevelRatios: []int{10, 10, 10}, // Each level is 10x larger
	}
}

func (l *LeveledStrategy) Name() string {
	return "Leveled"
}

func (l *LeveledStrategy) ShouldCompact(tables []*sstable.SSTable) bool {
	levels := l.groupByLevel(tables)
	
	for level, levelTables := range levels {
		if level == 0 {
			// Level 0 can have overlapping ranges, compact when too many
			if len(levelTables) >= 4 {
				return true
			}
		} else {
			// Other levels compact based on total size
			totalSize := int64(0)
			for _, table := range levelTables {
				totalSize += l.getTableSize(table)
			}
			
			maxSize := l.getLevelMaxSize(level)
			if totalSize > maxSize {
				return true
			}
		}
	}
	
	return false
}

func (l *LeveledStrategy) SelectTables(tables []*sstable.SSTable) []*sstable.SSTable {
	levels := l.groupByLevel(tables)
	
	// Find the level that needs compaction most urgently
	for level := 0; level <= l.MaxLevel; level++ {
		levelTables := levels[level]
		if len(levelTables) == 0 {
			continue
		}
		
		if level == 0 && len(levelTables) >= 4 {
			return levelTables
		}
		
		totalSize := int64(0)
		for _, table := range levelTables {
			totalSize += l.getTableSize(table)
		}
		
		if totalSize > l.getLevelMaxSize(level) {
			return levelTables
		}
	}
	
	return nil
}

func (l *LeveledStrategy) groupByLevel(tables []*sstable.SSTable) map[int][]*sstable.SSTable {
	levels := make(map[int][]*sstable.SSTable)
	
	for _, table := range tables {
		level := l.getTableLevel(table)
		levels[level] = append(levels[level], table)
	}
	
	return levels
}

func (l *LeveledStrategy) getTableLevel(table *sstable.SSTable) int {
	// Extract level from filename (e.g., "ss-L1-001.sst")
	filename := filepath.Base(table.Path)
	if strings.Contains(filename, "-L") {
		parts := strings.Split(filename, "-")
		for _, part := range parts {
			if strings.HasPrefix(part, "L") && len(part) > 1 {
				level := 0
				fmt.Sscanf(part[1:], "%d", &level)
				return level
			}
		}
	}
	return 0 // Default to level 0
}

func (l *LeveledStrategy) getLevelMaxSize(level int) int64 {
	if level < len(l.LevelSizes) {
		return l.LevelSizes[level]
	}
	
	// Calculate size for higher levels
	baseSize := l.LevelSizes[len(l.LevelSizes)-1]
	ratio := l.LevelRatios[len(l.LevelRatios)-1]
	
	for i := len(l.LevelSizes); i <= level; i++ {
		baseSize *= int64(ratio)
	}
	
	return baseSize
}

func (l *LeveledStrategy) getTableSize(table *sstable.SSTable) int64 {
	if info, err := os.Stat(table.Path); err == nil {
		return info.Size()
	}
	return 0
}

// TimeBasedStrategy compacts based on table age
type TimeBasedStrategy struct {
	MaxAge       int64 // Maximum age in seconds
	MinTables    int   // Minimum tables to compact
}

func NewTimeBasedStrategy() *TimeBasedStrategy {
	return &TimeBasedStrategy{
		MaxAge:    3600, // 1 hour
		MinTables: 3,
	}
}

func (t *TimeBasedStrategy) Name() string {
	return "Time-Based"
}

func (t *TimeBasedStrategy) ShouldCompact(tables []*sstable.SSTable) bool {
	if len(tables) < t.MinTables {
		return false
	}
	
	oldTables := 0
	for _, table := range tables {
		if t.isOld(table) {
			oldTables++
		}
	}
	
	return oldTables >= t.MinTables
}

func (t *TimeBasedStrategy) SelectTables(tables []*sstable.SSTable) []*sstable.SSTable {
	var oldTables []*sstable.SSTable
	
	for _, table := range tables {
		if t.isOld(table) {
			oldTables = append(oldTables, table)
		}
	}
	
	return oldTables
}

func (t *TimeBasedStrategy) isOld(table *sstable.SSTable) bool {
	if info, err := os.Stat(table.Path); err == nil {
		age := int64(info.ModTime().Unix())
		return age > t.MaxAge
	}
	return false
}
