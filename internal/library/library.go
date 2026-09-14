package library

import (
	"os"
	"path/filepath"
)

// This is the library package.
// The library exists of 4 folders: originals, thumbs (thumbnails), tmp (temporary), db (database).
// The library can be initialized using the New() function, it accepts a path parameter and build the lib relative to that path.

type Library struct {
	root string
}

func (l *Library) OriginalDir() string { return filepath.Join(l.root, "originals") }
func (l *Library) ThumbsDir() string   { return filepath.Join(l.root, "thumbnails") }
func (l *Library) TmpDir() string      { return filepath.Join(l.root, "tmp") }
func (l *Library) DbPath() string      { return filepath.Join(l.root, "library.db") }

func New(root string) (*Library, error) {
	l := &Library{root: root}

	for _, dir := range []string{l.OriginalDir(), l.ThumbsDir(), l.TmpDir(), l.DbPath()} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	return l, nil
}
