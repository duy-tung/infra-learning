package sstable

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"lsm/bloom"
	"lsm/memtable"
)

// SSTable represents an immutable sorted table on disk.
type SSTable struct {
	Path  string
	Bloom *bloom.Bloom
}

// New writes kvs to path and builds a Bloom filter.
func New(path string, kvs []memtable.KV) (*SSTable, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	for _, kv := range kvs {
		line := fmt.Sprintf("%s\t%s\n", kv.Key, kv.Value)
		if _, err := bw.WriteString(line); err != nil {
			return nil, err
		}
	}
	if err := bw.Flush(); err != nil {
		return nil, err
	}

	bl := bloom.New(uint(len(kvs)*8+1), 3)
	for _, kv := range kvs {
		bl.Add(kv.Key)
	}
	return &SSTable{Path: path, Bloom: bl}, nil
}

// Load opens an existing table and rebuilds its Bloom filter.
func Load(path string) (*SSTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var count int
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	bl := bloom.New(uint(count*8+1), 3)
	// second pass to add keys
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}
	scanner = bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "\t", 2)
		if len(parts) >= 1 {
			bl.Add(parts[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &SSTable{Path: path, Bloom: bl}, nil
}

// Get searches for key in the SSTable.
func (s *SSTable) Get(key string) (string, bool, error) {
	if !s.Bloom.Contains(key) {
		return "", false, nil
	}
	f, err := os.Open(s.Path)
	if err != nil {
		return "", false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "\t", 2)
		if len(parts) == 2 && parts[0] == key {
			return parts[1], true, nil
		}
	}
	return "", false, scanner.Err()
}
