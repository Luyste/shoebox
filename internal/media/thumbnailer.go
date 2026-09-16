package media

type Thumbnailer interface {
	Thumbnail(srcPath string, destPath string) error
}
