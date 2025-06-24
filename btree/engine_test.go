package btree

import (
	"os"
	"testing"
)

func TestEngineCRUD(t *testing.T) {
	path := "engine_test.gob"
	os.Remove(path)
	eng, err := Open(path, 3)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := eng.Put("a", "1"); err != nil {
		t.Fatalf("put: %v", err)
	}
	if val, ok, _ := eng.Get("a"); !ok || val != "1" {
		t.Fatalf("get failed: %v %v", val, ok)
	}
	if err := eng.Delete("a"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, ok, _ := eng.Get("a"); ok {
		t.Fatalf("should be deleted")
	}
	os.Remove(path)
}
