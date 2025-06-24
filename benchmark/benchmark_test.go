package benchmark

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"btree"
	"lsm/lsmtree"
)

type kv struct {
	k string
	v string
}

func genData(n int) []kv {
	r := rand.New(rand.NewSource(1))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	data := make([]kv, n)
	for i := 0; i < n; i++ {
		b := make([]rune, 16)
		for j := range b {
			b[j] = letters[r.Intn(len(letters))]
		}
		key := fmt.Sprintf("key%06d_%s", i, string(b))
		val := fmt.Sprintf("val_%s", string(b))
		data[i] = kv{k: key, v: val}
	}
	return data
}

func BenchmarkWriteLSM(b *testing.B) {
	data := genData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.RemoveAll("bench_lsm")
		tree, err := lsmtree.New("bench_lsm", 1000)
		if err != nil {
			b.Fatalf("new lsm: %v", err)
		}
		for _, kv := range data {
			if err := tree.Put(kv.k, kv.v); err != nil {
				b.Fatalf("put: %v", err)
			}
		}
	}
}

func BenchmarkWriteBTree(b *testing.B) {
	data := genData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.Remove("btree.gob")
		eng, err := btree.Open("btree.gob", 3)
		if err != nil {
			b.Fatalf("open: %v", err)
		}
		for _, kv := range data {
			if err := eng.Put(kv.k, kv.v); err != nil {
				b.Fatalf("put: %v", err)
			}
		}
	}
}

func prepLSM(data []kv) (*lsmtree.LSMTree, error) {
	os.RemoveAll("bench_lsm")
	t, err := lsmtree.New("bench_lsm", 1000)
	if err != nil {
		return nil, err
	}
	for _, kv := range data {
		if err := t.Put(kv.k, kv.v); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func prepBTree(data []kv) (*btree.Engine, error) {
	os.Remove("btree.gob")
	e, err := btree.Open("btree.gob", 3)
	if err != nil {
		return nil, err
	}
	for _, kv := range data {
		if err := e.Put(kv.k, kv.v); err != nil {
			return nil, err
		}
	}
	return e, nil
}

func BenchmarkReadLSM(b *testing.B) {
	data := genData(10000)
	tree, err := prepLSM(data)
	if err != nil {
		b.Fatalf("prep lsm: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, kv := range data {
			if _, ok, err := tree.Get(kv.k); err != nil || !ok {
				b.Fatalf("get: %v %v", err, ok)
			}
		}
	}
}

func BenchmarkReadBTree(b *testing.B) {
	data := genData(10000)
	eng, err := prepBTree(data)
	if err != nil {
		b.Fatalf("prep btree: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, kv := range data {
			if _, ok, err := eng.Get(kv.k); err != nil || !ok {
				b.Fatalf("get: %v %v", err, ok)
			}
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
