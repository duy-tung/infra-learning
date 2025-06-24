package btree

import (
	"encoding/gob"
	"os"
)

// Node represents a single B-tree node persisted as a page.
type Node struct {
	Leaf     bool
	Keys     []string
	Values   []string
	Children []*Node
}

// BTree implements a simple B-tree with disk persistence.
type BTree struct {
	Order int
	Root  *Node
}

// New creates an empty B-tree of given order.
func New(order int) *BTree {
	if order < 2 {
		panic("order must be >= 2")
	}
	return &BTree{Order: order, Root: &Node{Leaf: true}}
}

// Save writes the entire tree to file using gob encoding.
func (t *BTree) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(t)
}

// Load reads the tree state from file.
func Load(path string) (*BTree, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var bt BTree
	if err := gob.NewDecoder(f).Decode(&bt); err != nil {
		return nil, err
	}
	return &bt, nil
}

// Search retrieves a value by key.
func (t *BTree) Search(key string) (string, bool) {
	n := t.Root
	for {
		i := 0
		for i < len(n.Keys) && key > n.Keys[i] {
			i++
		}
		if i < len(n.Keys) && key == n.Keys[i] {
			return n.Values[i], true
		}
		if n.Leaf {
			return "", false
		}
		n = n.Children[i]
	}
}

// splitChild splits child y of x at index i.
func splitChild(x *Node, i int, t int) {
	y := x.Children[i]
	z := &Node{Leaf: y.Leaf}
	mid := t - 1

	midKey := y.Keys[mid]
	midVal := y.Values[mid]
	z.Keys = append(z.Keys, y.Keys[mid+1:]...)
	z.Values = append(z.Values, y.Values[mid+1:]...)
	y.Keys = y.Keys[:mid]
	y.Values = y.Values[:mid]

	if !y.Leaf {
		z.Children = append(z.Children, y.Children[mid+1:]...)
		y.Children = y.Children[:mid+1]
	}

	x.Children = append(x.Children[:i+1], append([]*Node{z}, x.Children[i+1:]...)...)
	x.Keys = append(x.Keys[:i], append([]string{midKey}, x.Keys[i:]...)...)
	x.Values = append(x.Values[:i], append([]string{midVal}, x.Values[i:]...)...)
}

// insertNonFull inserts key/value into node assumed not full.
func insertNonFull(n *Node, key, value string, t int) {
	i := len(n.Keys) - 1
	if n.Leaf {
		n.Keys = append(n.Keys, "")
		n.Values = append(n.Values, "")
		for i >= 0 && key < n.Keys[i] {
			n.Keys[i+1] = n.Keys[i]
			n.Values[i+1] = n.Values[i]
			i--
		}
		n.Keys[i+1] = key
		n.Values[i+1] = value
		return
	}
	for i >= 0 && key < n.Keys[i] {
		i--
	}
	i++
	if len(n.Children[i].Keys) == 2*t-1 {
		splitChild(n, i, t)
		if key > n.Keys[i] {
			i++
		}
	}
	insertNonFull(n.Children[i], key, value, t)
}

// Insert adds a key/value to the tree.
func (t *BTree) Insert(key, value string) {
	r := t.Root
	if len(r.Keys) == 2*t.Order-1 {
		s := &Node{}
		t.Root = s
		s.Children = []*Node{r}
		splitChild(s, 0, t.Order)
		insertNonFull(s, key, value, t.Order)
	} else {
		insertNonFull(r, key, value, t.Order)
	}
}

// mergeChildren merges child i and i+1 of node x.
func mergeChildren(x *Node, i int, t int) {
	y := x.Children[i]
	z := x.Children[i+1]

	y.Keys = append(y.Keys, x.Keys[i])
	y.Values = append(y.Values, x.Values[i])
	y.Keys = append(y.Keys, z.Keys...)
	y.Values = append(y.Values, z.Values...)
	if !y.Leaf {
		y.Children = append(y.Children, z.Children...)
	}

	x.Keys = append(x.Keys[:i], x.Keys[i+1:]...)
	x.Values = append(x.Values[:i], x.Values[i+1:]...)
	x.Children = append(x.Children[:i+1], x.Children[i+2:]...)
}

// borrowFromPrev borrows a key from child i-1 of x to child i.
func borrowFromPrev(x *Node, i int) {
	child := x.Children[i]
	sibling := x.Children[i-1]

	child.Keys = append([]string{x.Keys[i-1]}, child.Keys...)
	child.Values = append([]string{x.Values[i-1]}, child.Values...)
	if !child.Leaf {
		child.Children = append([]*Node{sibling.Children[len(sibling.Children)-1]}, child.Children...)
		sibling.Children = sibling.Children[:len(sibling.Children)-1]
	}
	x.Keys[i-1] = sibling.Keys[len(sibling.Keys)-1]
	x.Values[i-1] = sibling.Values[len(sibling.Values)-1]
	sibling.Keys = sibling.Keys[:len(sibling.Keys)-1]
	sibling.Values = sibling.Values[:len(sibling.Values)-1]
}

// borrowFromNext borrows a key from child i+1 of x to child i.
func borrowFromNext(x *Node, i int) {
	child := x.Children[i]
	sibling := x.Children[i+1]

	child.Keys = append(child.Keys, x.Keys[i])
	child.Values = append(child.Values, x.Values[i])
	if !child.Leaf {
		child.Children = append(child.Children, sibling.Children[0])
		sibling.Children = sibling.Children[1:]
	}
	x.Keys[i] = sibling.Keys[0]
	x.Values[i] = sibling.Values[0]
	sibling.Keys = sibling.Keys[1:]
	sibling.Values = sibling.Values[1:]
}

// removeFromNode removes key from subtree rooted with node.
func removeFromNode(n *Node, key string, t int) {
	idx := 0
	for idx < len(n.Keys) && key > n.Keys[idx] {
		idx++
	}

	if idx < len(n.Keys) && n.Keys[idx] == key {
		if n.Leaf {
			n.Keys = append(n.Keys[:idx], n.Keys[idx+1:]...)
			n.Values = append(n.Values[:idx], n.Values[idx+1:]...)
			return
		}
		if len(n.Children[idx].Keys) >= t {
			pred := n.Children[idx]
			for !pred.Leaf {
				pred = pred.Children[len(pred.Children)-1]
			}
			pk := pred.Keys[len(pred.Keys)-1]
			pv := pred.Values[len(pred.Values)-1]
			removeFromNode(n.Children[idx], pk, t)
			n.Keys[idx] = pk
			n.Values[idx] = pv
			return
		}
		if len(n.Children[idx+1].Keys) >= t {
			succ := n.Children[idx+1]
			for !succ.Leaf {
				succ = succ.Children[0]
			}
			sk := succ.Keys[0]
			sv := succ.Values[0]
			removeFromNode(n.Children[idx+1], sk, t)
			n.Keys[idx] = sk
			n.Values[idx] = sv
			return
		}
		mergeChildren(n, idx, t)
		removeFromNode(n.Children[idx], key, t)
		return
	}

	if n.Leaf {
		return
	}

	flag := idx == len(n.Keys)

	if len(n.Children[idx].Keys) < t {
		if idx != 0 && len(n.Children[idx-1].Keys) >= t {
			borrowFromPrev(n, idx)
		} else if idx != len(n.Keys) && len(n.Children[idx+1].Keys) >= t {
			borrowFromNext(n, idx)
		} else {
			if idx != len(n.Keys) {
				mergeChildren(n, idx, t)
			} else {
				mergeChildren(n, idx-1, t)
			}
		}
	}
	if flag && idx > len(n.Keys) {
		removeFromNode(n.Children[idx-1], key, t)
	} else {
		removeFromNode(n.Children[idx], key, t)
	}
}

// Delete removes key from the tree.
func (t *BTree) Delete(key string) {
	if t.Root == nil {
		return
	}
	removeFromNode(t.Root, key, t.Order)
	if len(t.Root.Keys) == 0 && !t.Root.Leaf {
		t.Root = t.Root.Children[0]
	}
}
