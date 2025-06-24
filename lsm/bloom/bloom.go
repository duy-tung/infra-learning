package bloom

import (
	"hash/fnv"
)

// Bloom is a simple Bloom filter.
type Bloom struct {
	bits []byte
	k    uint
}

// New creates a Bloom filter with size m bits and k hash functions.
func New(m uint, k uint) *Bloom {
	return &Bloom{bits: make([]byte, m), k: k}
}

// Add inserts a string into the filter.
func (b *Bloom) Add(s string) {
	indexes := b.hashes(s)
	for _, idx := range indexes {
		b.bits[idx%uint(len(b.bits))] = 1
	}
}

// Contains checks if a string is possibly in the set.
func (b *Bloom) Contains(s string) bool {
	indexes := b.hashes(s)
	for _, idx := range indexes {
		if b.bits[idx%uint(len(b.bits))] == 0 {
			return false
		}
	}
	return true
}

func (b *Bloom) hashes(s string) []uint {
	hashes := make([]uint, b.k)
	h := fnv.New64a()
	for i := uint(0); i < b.k; i++ {
		h.Reset()
		h.Write([]byte{byte(i)})
		h.Write([]byte(s))
		hashes[i] = uint(h.Sum64())
	}
	return hashes
}
