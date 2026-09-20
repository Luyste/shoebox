package web

import (
	"io/fs"
	"log"
	"net/http"
)

func (w *Web) Router() *http.ServeMux {
	mu := http.NewServeMux()

	sub, err := fs.Sub(w.wfs, "static")

	if err != nil {
		log.Fatal("unable to init static assets")
	}

	mu.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(sub)))

	mu.HandleFunc("GET /", w.Home)

	return mu
}
