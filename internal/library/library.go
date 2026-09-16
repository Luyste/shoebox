package library

import (
	"os"
	"path/filepath"
)

type Library struct {
	root string
}

func (l *Library) OriginalDir() string { return filepath.Join(l.root, "originals") }
func (l *Library) ThumbsDir() string   { return filepath.Join(l.root, "thumbnails") }
func (l *Library) TmpDir() string      { return filepath.Join(l.root, "tmp") }
func (l *Library) DbPath() string      { return filepath.Join(l.root, "library.db") }

func New(root string) (*Library, error) {
	l := &Library{root: root}

	for _, dir := range []string{l.OriginalDir(), l.ThumbsDir(), l.TmpDir()} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	return l, nil
}
