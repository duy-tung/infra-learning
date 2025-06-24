package btree

import (
	"os"
)

// Engine wraps a B-tree providing persistent operations similar to lsmtree.
type Engine struct {
	tree *BTree
	path string
}

// Open creates or loads a B-tree at the given file path.
func Open(path string, order int) (*Engine, error) {
	var t *BTree
	if _, err := os.Stat(path); err == nil {
		bt, err := Load(path)
		if err != nil {
			return nil, err
		}
		t = bt
	} else if os.IsNotExist(err) {
		t = New(order)
	} else {
		return nil, err
	}
	return &Engine{tree: t, path: path}, nil
}

// Put inserts or updates a key/value pair and persists the tree.
func (e *Engine) Put(key, value string) error {
	e.tree.Insert(key, value)
	return e.tree.Save(e.path)
}

// Get retrieves a value by key.
func (e *Engine) Get(key string) (string, bool, error) {
	val, ok := e.tree.Search(key)
	return val, ok, nil
}

// Delete removes a key from the tree and persists the change.
func (e *Engine) Delete(key string) error {
	e.tree.Delete(key)
	return e.tree.Save(e.path)
}
