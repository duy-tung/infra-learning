package bloom

import (
	"fmt"
	"hash/fnv"
	"math"
)

// EnhancedBloom provides additional functionality over the basic Bloom filter
type EnhancedBloom struct {
	*Bloom
	expectedElements uint
	actualElements   uint
}

// NewEnhanced creates an optimally-sized Bloom filter for expected elements and desired false positive rate
func NewEnhanced(expectedElements uint, falsePositiveRate float64) *EnhancedBloom {
	// Calculate optimal size: m = -n * ln(p) / (ln(2)^2)
	m := uint(-float64(expectedElements) * math.Log(falsePositiveRate) / (math.Log(2) * math.Log(2)))
	
	// Calculate optimal number of hash functions: k = (m/n) * ln(2)
	k := uint(float64(m) / float64(expectedElements) * math.Log(2))
	
	// Ensure minimum values
	if m < 8 {
		m = 8
	}
	if k < 1 {
		k = 1
	}
	if k > 10 {
		k = 10 // Practical limit
	}
	
	return &EnhancedBloom{
		Bloom:            New(m, k),
		expectedElements: expectedElements,
		actualElements:   0,
	}
}

// Add inserts an element and tracks count
func (eb *EnhancedBloom) Add(s string) {
	eb.Bloom.Add(s)
	eb.actualElements++
}

// Stats returns statistics about the filter
func (eb *EnhancedBloom) Stats() BloomStats {
	// Calculate current false positive probability
	// p = (1 - e^(-k*n/m))^k
	k := float64(eb.k)
	n := float64(eb.actualElements)
	m := float64(len(eb.bits))
	
	falsePositiveRate := math.Pow(1-math.Exp(-k*n/m), k)
	
	// Calculate fill ratio
	setBits := 0
	for _, bit := range eb.bits {
		if bit == 1 {
			setBits++
		}
	}
	fillRatio := float64(setBits) / float64(len(eb.bits))
	
	return BloomStats{
		Size:              uint(len(eb.bits)),
		HashFunctions:     eb.k,
		ExpectedElements:  eb.expectedElements,
		ActualElements:    eb.actualElements,
		FalsePositiveRate: falsePositiveRate,
		FillRatio:         fillRatio,
		SetBits:           uint(setBits),
	}
}

// BloomStats contains statistics about a Bloom filter
type BloomStats struct {
	Size              uint
	HashFunctions     uint
	ExpectedElements  uint
	ActualElements    uint
	FalsePositiveRate float64
	FillRatio         float64
	SetBits           uint
}

// String returns a formatted string representation of the stats
func (bs BloomStats) String() string {
	return fmt.Sprintf(`Bloom Filter Statistics:
  Size: %d bits
  Hash Functions: %d
  Expected Elements: %d
  Actual Elements: %d
  False Positive Rate: %.4f (%.2f%%)
  Fill Ratio: %.4f (%.2f%%)
  Set Bits: %d/%d`,
		bs.Size, bs.HashFunctions, bs.ExpectedElements, bs.ActualElements,
		bs.FalsePositiveRate, bs.FalsePositiveRate*100,
		bs.FillRatio, bs.FillRatio*100,
		bs.SetBits, bs.Size)
}

// CountingBloom implements a counting Bloom filter that supports deletions
type CountingBloom struct {
	counters []uint8
	k        uint
}

// NewCounting creates a counting Bloom filter
func NewCounting(m uint, k uint) *CountingBloom {
	return &CountingBloom{
		counters: make([]uint8, m),
		k:        k,
	}
}

// Add increments counters for the element
func (cb *CountingBloom) Add(s string) {
	indexes := cb.hashes(s)
	for _, idx := range indexes {
		pos := idx % uint(len(cb.counters))
		if cb.counters[pos] < 255 { // Prevent overflow
			cb.counters[pos]++
		}
	}
}

// Remove decrements counters for the element
func (cb *CountingBloom) Remove(s string) {
	indexes := cb.hashes(s)
	for _, idx := range indexes {
		pos := idx % uint(len(cb.counters))
		if cb.counters[pos] > 0 {
			cb.counters[pos]--
		}
	}
}

// Contains checks if element might be in the set
func (cb *CountingBloom) Contains(s string) bool {
	indexes := cb.hashes(s)
	for _, idx := range indexes {
		if cb.counters[idx%uint(len(cb.counters))] == 0 {
			return false
		}
	}
	return true
}

func (cb *CountingBloom) hashes(s string) []uint {
	hashes := make([]uint, cb.k)
	h := fnv.New64a()
	for i := uint(0); i < cb.k; i++ {
		h.Reset()
		h.Write([]byte{byte(i)})
		h.Write([]byte(s))
		hashes[i] = uint(h.Sum64())
	}
	return hashes
}
