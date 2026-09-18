package indexer

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/luyste/shoebox/internal/library"
	"github.com/luyste/shoebox/internal/media"
)

type Indexer struct {
	db          *sql.DB
	lib         *library.Library
	thumbnailer media.Thumbnailer
	mu          sync.Mutex
}

type FileResult struct {
	Path string
	Id   string
	Err  error
}

type Result struct {
	Indexed int
	Skipped int
	Failed  int
	Files   []FileResult
	Err     error
}

var ErrDuplicate = errors.New("file already indexed")

func New(db *sql.DB, lib *library.Library, thumbnailer media.Thumbnailer) *Indexer {
	return &Indexer{
		db:          db,
		lib:         lib,
		thumbnailer: thumbnailer,
	}
}

func (idx *Indexer) Index(root string) Result {
	paths := make(chan string)
	results := make(chan FileResult)
	walkErr := make(chan error, 1)

	const workerCount int = 4
	var wg sync.WaitGroup
	var report []FileResult
	var indexedCount, skippedCount, failedCount int

	go func() {
		walkErr <- walkFileTree(root, paths)
	}()

	for range workerCount {
		wg.Go(func() {
			for p := range paths {
				id, err := idx.IndexFile(p)

				res := FileResult{
					Path: p,
					Id:   id,
					Err:  err,
				}
				results <- res
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		switch {
		case errors.Is(res.Err, ErrDuplicate):
			skippedCount++
		case res.Err == nil:
			indexedCount++
		default:
			failedCount++
		}

		report = append(report, res)
	}

	if err := <-walkErr; err != nil {
		return Result{
			Indexed: indexedCount,
			Skipped: skippedCount,
			Failed:  failedCount,
			Files:   nil,
			Err:     err,
		}
	}

	return Result{
		Files:   report,
		Indexed: indexedCount,
		Skipped: skippedCount,
		Failed:  failedCount,
		Err:     nil,
	}
}

func (idx *Indexer) cleanUp(dest string, id string) error {
	if err := os.Remove(dest); err != nil {
		return fmt.Errorf("removing file: %w", err)
	}

	_, err := idx.db.Exec("DELETE FROM media WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting orphaned row %v: %w", id, err)
	}
	return nil
}

func (idx *Indexer) IndexFile(srcPath string) (string, error) {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("something went wrong reading file: %w", err)
	}

	info, err := media.Inspect(content)
	if err != nil {
		return "", fmt.Errorf("something went wrong inspecting file: %w", err)
	}

	var existingID string

	id, err := func() (string, error) {
		idx.mu.Lock()
		defer idx.mu.Unlock()

		err = idx.db.QueryRow("SELECT id FROM media WHERE sha256 = ?", info.Checksum).Scan(&existingID)

		if errors.Is(err, sql.ErrNoRows) {
			newID := uuid.NewString()
			originalName := filepath.Base(srcPath)
			ogfp := filepath.Join(idx.lib.OriginalDir(), newID+".jpg")

			_, err := idx.db.Exec("INSERT INTO media (id, sha256, kind, original_name, original_path, size_bytes, mime_type) VALUES (?, ?, ?, ?, ?, ?, ?)", newID, info.Checksum, info.Kind, originalName, ogfp, info.Size, info.MimeType)
			if err != nil {
				return "", err
			}

			return newID, nil
		}

		if err == nil {
			return existingID, ErrDuplicate
		}

		return "", fmt.Errorf("checking for duplicate: %w", err)
	}()

	if err == nil {
		ogfp := filepath.Join(idx.lib.OriginalDir(), id+".jpg")

		if err := os.WriteFile(ogfp, content, 0o644); err != nil {
			if err := idx.cleanUp(ogfp, id); err != nil {
				return "", fmt.Errorf("cleanup file: %w", err)
			}
			return "", fmt.Errorf("writing file: %w", err)
		}

		thfp := filepath.Join(idx.lib.ThumbsDir(), id+".jpg")
		if err := idx.thumbnailer.Thumbnail(ogfp, thfp); err != nil {
			if err := idx.cleanUp(thfp, id); err != nil {
				return "", fmt.Errorf("cleanup file: %w", err)
			}
			return "", fmt.Errorf("writing thumbnail: %w", err)
		}

		_, err = idx.db.Exec("UPDATE media SET thumb_state = 'done' WHERE id = ?", id)
		if err != nil {
			return "", fmt.Errorf("update database: %w", err)
		}

		return id, nil
	}

	return existingID, err
}

func walkFileTree(root string, pathChan chan string) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking file tree: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		pathChan <- path
		return nil
	})

	close(pathChan)
	return err
}
