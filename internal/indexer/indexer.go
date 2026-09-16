package indexer

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/luyste/shoebox/internal/library"
	"github.com/luyste/shoebox/internal/media"
)

type Indexer struct {
	db          *sql.DB
	lib         *library.Library
	thumbnailer media.Thumbnailer
}

var ErrDuplicate = errors.New("file already indexed")

func New(db *sql.DB, lib *library.Library, thumbnailer media.Thumbnailer) *Indexer {
	return &Indexer{
		db:          db,
		lib:         lib,
		thumbnailer: thumbnailer,
	}
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
	err = idx.db.QueryRow("SELECT id FROM media WHERE sha256 = ?", info.Checksum).Scan(&existingID)

	switch {
	case err == nil:
		return existingID, ErrDuplicate
	case errors.Is(err, sql.ErrNoRows):
		newID := uuid.NewString()
		ogfp := filepath.Join(idx.lib.OriginalDir(), newID+".jpg")
		if err := os.WriteFile(ogfp, content, 0o644); err != nil {
			return "", fmt.Errorf("writing file: %w", err)
		}

		thfp := filepath.Join(idx.lib.ThumbsDir(), newID+".jpg")
		if err := idx.thumbnailer.Thumbnail(ogfp, thfp); err != nil {
			return "", fmt.Errorf("writing thumbnail: %w", err)
		}

		originalName := filepath.Base(srcPath)

		_, err := idx.db.Exec("INSERT INTO media (id, sha256, kind, original_name, original_path, size_bytes, mime_type, thumb_state) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", newID, info.Checksum, info.Kind, originalName, ogfp, info.Size, info.MimeType, "done")
		if err != nil {
			return "", fmt.Errorf("insert database: %w", err)
		}

		return newID, nil

	default:
		return "", fmt.Errorf("checking for duplicate: %w", err)
	}
}
