package indexer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkFileTreeScratch(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	os.WriteFile(filepath.Join(root, "a.jpg"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(sub, "b.jpg"), []byte("b"), 0o644)

	paths := make(chan string)
	go walkFileTree(root, paths)

	var found []string
	for p := range paths {
		found = append(found, p)
	}

	t.Logf("found %d paths: %v", len(found), found)
	if len(found) != 2 {
		t.Errorf("want 2 files, got %d", len(found))
	}
}
