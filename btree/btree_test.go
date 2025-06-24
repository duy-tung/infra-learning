package btree

import (
	"fmt"
	"os"
	"testing"
)

func TestInsertSearchDelete(t *testing.T) {
	bt := New(3)
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("k%03d", i)
		v := fmt.Sprintf("v%03d", i)
		bt.Insert(k, v)
	}
	val, ok := bt.Search("k050")
	if !ok || val != "v050" {
		t.Fatalf("expected v050, got %s", val)
	}
	bt.Delete("k050")
	if _, ok := bt.Search("k050"); ok {
		t.Fatalf("key should be deleted")
	}
	tmp := "btree_test.gob"
	if err := bt.Save(tmp); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Fatalf("file not written: %v", err)
	}
	loaded, err := Load(tmp)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if _, ok := loaded.Search("k051"); !ok {
		t.Fatalf("loaded tree missing key")
	}
	os.Remove(tmp)
}
