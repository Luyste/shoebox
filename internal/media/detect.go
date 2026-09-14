package media

import "net/http"

var AllowedMimeTypes = map[string]string{"image/jpeg": "photo", "video/mp4": "video"}

func detectMimeType(byte []byte) string {
	return http.DetectContentType(byte)
}

func kindFromMimeType(mimeType string) (string, error) {

}
