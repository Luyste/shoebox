package web

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/luyste/shoebox/internal/library"
)

type Web struct {
	db  *sql.DB
	lib *library.Library
	wfs *embed.FS
}

func New(db *sql.DB, lib *library.Library, wfs *embed.FS) *Web {
	return &Web{
		db:  db,
		lib: lib,
		wfs: wfs,
	}
}

func (w *Web) Router() *http.ServeMux {
	mu := http.NewServeMux()

	sub, err := fs.Sub(w.wfs, "static")

	if err != nil {
		log.Fatal("unable to init static assets")
	}

	mu.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(sub)))

	mu.HandleFunc("GET /", func(wr http.ResponseWriter, r *http.Request) {
		response := "ok"
		wr.Write([]byte(response))
	})

	return mu
}
