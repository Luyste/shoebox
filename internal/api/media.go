package api

import (
	"net/http"
)

type newMedia struct {
	id           string
	sha256       string
	kind         string
	originalName string
	size_bytes   int64
	mime_type    string
}

func (a *API) UploadMedia(w http.ResponseWriter, r *http.Request) {

}

func (a *API) DownloadMedia(w http.ResponseWriter, r *http.Request) {

}
