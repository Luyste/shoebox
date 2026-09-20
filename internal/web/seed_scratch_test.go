package web

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/luyste/shoebox/internal/db"
	"github.com/luyste/shoebox/internal/indexer"
	"github.com/luyste/shoebox/internal/library"
	"github.com/luyste/shoebox/internal/media"
)

// TestSeedScratch is a manual helper: seeds a real library at a fixed path
// so the server can be run against it by hand for end-to-end checks.
func TestSeedScratch(t *testing.T) {
	libRoot := "/tmp/shoebox-seed"
	os.RemoveAll(libRoot)
	lib, err := library.New(libRoot)
	if err != nil {
		t.Fatalf("library.New: %v", err)
	}

	conn, err := db.Connect(lib.DbPath())
	if err != nil {
		t.Fatalf("db.Connect: %v", err)
	}
	defer conn.Close()

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	importDir := filepath.Join(libRoot, "import")
	os.MkdirAll(importDir, 0o755)
	for i := 0; i < 3; i++ {
		p := filepath.Join(importDir, string(rune('a'+i))+".jpg")
		f, err := os.Create(p)
		if err != nil {
			t.Fatalf("creating test image: %v", err)
		}
		target := image.NewRGBA(image.Rect(0, 0, 10+i, 10))
		if err := jpeg.Encode(f, target, nil); err != nil {
			t.Fatalf("encoding test image: %v", err)
		}
		f.Close()
	}

	idx := indexer.New(conn, lib, media.VipsThumbnailer{})
	res := idx.Index(importDir)
	t.Logf("seed result: %+v", res)
}
