package api

import (
	"github.com/go-chi/chi/v5"
)

func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/photo/{id}", a.DownloadMedia)
	r.Post("/photo", a.UploadMedia)

	return r
}
