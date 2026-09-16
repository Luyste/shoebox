package indexer

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/luyste/shoebox/internal/db"
	"github.com/luyste/shoebox/internal/library"
	"github.com/luyste/shoebox/internal/media"
)

func TestIndexScratch(t *testing.T) {
	libRoot := t.TempDir()
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

	// folder of source images to import, separate from the library itself
	importDir := t.TempDir()
	for i := 0; i < 5; i++ {
		p := filepath.Join(importDir, "img"+string(rune('a'+i))+".jpg")
		f, err := os.Create(p)
		if err != nil {
			t.Fatalf("creating test image: %v", err)
		}
		target := image.NewRGBA(image.Rect(0, 0, 10, 10))
		if err := jpeg.Encode(f, target, nil); err != nil {
			t.Fatalf("encoding test image: %v", err)
		}
		f.Close()
	}

	idx := New(conn, lib, media.VipsThumbnailer{})

	done := make(chan struct{})
	go func() {
		idx.Index(importDir)
		close(done)
	}()

	select {
	case <-done:
		t.Log("Index completed")
	case <-time.After(10 * time.Second):
		t.Fatal("Index did not complete within 10s — possible deadlock")
	}
}

func TestIndexScratchDBCheck(t *testing.T) {
	libRoot := t.TempDir()
	lib, _ := library.New(libRoot)
	conn, _ := db.Connect(lib.DbPath())
	defer conn.Close()
	db.Migrate(conn)

	importDir := t.TempDir()
	p := filepath.Join(importDir, "img.jpg")
	f, _ := os.Create(p)
	target := image.NewRGBA(image.Rect(0, 0, 10, 10))
	jpeg.Encode(f, target, nil)
	f.Close()

	idx := New(conn, lib, media.VipsThumbnailer{})
	id, err := idx.IndexFile(p)
	t.Logf("IndexFile returned id=%q err=%v", id, err)

	var count int
	conn.QueryRow("SELECT COUNT(*) FROM media").Scan(&count)
	t.Logf("rows in media table: %d", count)

	entries, _ := os.ReadDir(lib.OriginalDir())
	t.Logf("files in originals dir: %d", len(entries))
	for _, e := range entries {
		t.Logf("  - %s", e.Name())
	}
}

func TestIndexScratchThumbState(t *testing.T) {
	libRoot := t.TempDir()
	lib, _ := library.New(libRoot)
	conn, _ := db.Connect(lib.DbPath())
	defer conn.Close()
	db.Migrate(conn)

	importDir := t.TempDir()
	p := filepath.Join(importDir, "img.jpg")
	f, _ := os.Create(p)
	target := image.NewRGBA(image.Rect(0, 0, 10, 10))
	jpeg.Encode(f, target, nil)
	f.Close()

	idx := New(conn, lib, media.VipsThumbnailer{})
	id, err := idx.IndexFile(p)
	if err != nil {
		t.Fatalf("IndexFile: %v", err)
	}

	var state string
	if err := conn.QueryRow("SELECT thumb_state FROM media WHERE id = ?", id).Scan(&state); err != nil {
		t.Fatalf("querying thumb_state: %v", err)
	}
	t.Logf("thumb_state: %s", state)
	if state != "done" {
		t.Errorf("want thumb_state=done, got %s", state)
	}

	thumbs, _ := os.ReadDir(lib.ThumbsDir())
	t.Logf("files in thumbs dir: %d", len(thumbs))
	if len(thumbs) != 1 {
		t.Errorf("want 1 thumbnail file, got %d", len(thumbs))
	}
}
