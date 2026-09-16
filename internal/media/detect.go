package media

import (
	"fmt"
	"net/http"
)

type Info struct {
	Kind     string
	MimeType string
	Size     int
	Checksum string
}

var allowedMimeTypes = map[string]string{"image/jpeg": "photo", "video/mp4": "video"}

func detectMimeType(content []byte) string {
	return http.DetectContentType(content)
}

func kindFromMimeType(mimeType string) (string, error) {
	val, ok := allowedMimeTypes[mimeType]

	if !ok {
		return "", fmt.Errorf("unsupported mime type detected: %v", mimeType)
	}
	return val, nil
}

func Inspect(content []byte) (Info, error) {
	mimeType := detectMimeType(content)
	kind, err := kindFromMimeType(mimeType)
	if err != nil {
		return Info{}, err
	}
	size := len(content)
	checksum := buildChecksum(content)

	return Info{
		Kind:     kind,
		MimeType: mimeType,
		Size:     size,
		Checksum: checksum,
	}, nil
}
