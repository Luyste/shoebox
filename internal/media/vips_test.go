package media

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func TestVipsThumbnailer(t *testing.T) {
	srcPath := filepath.Join(t.TempDir(), "img.jpg")
	destPath := filepath.Join(t.TempDir(), "thumb.jpg")

	tests := []struct {
		name     string
		srcPath  string
		destPath string
		wantErr  bool
	}{
		{name: "succesfully generates thumbnail", srcPath: srcPath, destPath: destPath, wantErr: false},
		{name: "fails generating thumbnail", srcPath: "/", destPath: "/", wantErr: true},
	}

	target := image.NewRGBA(image.Rect(0, 0, 10, 10))

	f, err := os.Create(srcPath)
	if err != nil {
		t.Fatalf("creating source file: %v", err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, target, nil); err != nil {
		t.Fatalf("encoding image: %v", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := VipsThumbnailer.Thumbnail(VipsThumbnailer{}, tc.srcPath, tc.destPath); tc.wantErr != (err != nil) {
				t.Errorf("wantErr: %v, err: %v", tc.wantErr, err)
			}

			if !tc.wantErr {
				info, err := os.Stat(tc.destPath)

				if err != nil {
					t.Errorf("something went wrong fetching file info: %v", err)
				}

				if info.Name() != "thumb.jpg" {
					t.Errorf("filename expected to be thumb.jpg, got:%v", info.Name())
				}
			}
		})
	}
}
