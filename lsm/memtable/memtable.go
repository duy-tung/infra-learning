package memtable

import "sort"

// Memtable holds key-value pairs in memory until flush threshold.
type Memtable struct {
	Data           map[string]string
	FlushThreshold int
}

// New creates a new Memtable with given flush threshold.
func New(threshold int) *Memtable {
	return &Memtable{
		Data:           make(map[string]string),
		FlushThreshold: threshold,
	}
}

// Put inserts or updates a key-value pair.
func (m *Memtable) Put(key, value string) {
	m.Data[key] = value
}

// Get retrieves a value and boolean indicating presence.
func (m *Memtable) Get(key string) (string, bool) {
	v, ok := m.Data[key]
	return v, ok
}

// IsFull checks if memtable reached flush threshold.
func (m *Memtable) IsFull() bool {
	return len(m.Data) >= m.FlushThreshold
}

// Flush returns sorted contents and resets the memtable.
func (m *Memtable) Flush() []KV {
	kvs := make([]KV, 0, len(m.Data))
	for k, v := range m.Data {
		kvs = append(kvs, KV{Key: k, Value: v})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].Key < kvs[j].Key })
	m.Data = make(map[string]string)
	return kvs
}

// KV is a key-value pair.
type KV struct {
	Key   string
	Value string
}
