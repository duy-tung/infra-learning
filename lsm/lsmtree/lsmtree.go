package lsmtree

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lsm/memtable"
	"lsm/sstable"
)

// LSMTree coordinates memtable and SSTables.
type LSMTree struct {
	Mem    *memtable.Memtable
	Tables []*sstable.SSTable
	Dir    string
	nextID int
}

// New creates tree and loads existing tables from dir.
func New(dir string, threshold int) (*LSMTree, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	t := &LSMTree{Mem: memtable.New(threshold), Dir: dir}

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
	t.Mem.Put(key, value)
	if t.Mem.IsFull() {
		return t.flush()
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
	if v, ok := t.Mem.Get(key); ok {
		return v, true, nil
	}
	for i := len(t.Tables) - 1; i >= 0; i-- {
		if v, ok, err := t.Tables[i].Get(key); err != nil {
			return "", false, err
		} else if ok {
			return v, true, nil
		}
	}
	return "", false, nil
}

// Compact merges all tables into one.
func (t *LSMTree) Compact() error {
	if len(t.Tables) < 2 {
		return nil
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
