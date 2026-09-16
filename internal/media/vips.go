package media

import (
	"fmt"
	"os/exec"
)

type VipsThumbnailer struct{}

func (v VipsThumbnailer) Thumbnail(srcPath string, destPath string) error {
	cmd := exec.Command("vipsthumbnail", srcPath, "-s", "400", "-o", destPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generation of thumbnail failed for file %v: %w", srcPath, err)
	}
	return nil
}
