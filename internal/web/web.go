package web

import (
	"database/sql"
	"net/http"

	"github.com/luyste/shoebox/internal/library"
)

type Web struct {
	db  *sql.DB
	lib *library.Library
}

func New(db *sql.DB, lib *library.Library) *Web {
	return &Web{
		db:  db,
		lib: lib,
	}
}

func (w *Web) Router() *http.ServeMux {
	mu := http.NewServeMux()

	mu.HandleFunc("GET /", func(wr http.ResponseWriter, r *http.Request) {
		response := "ok"
		wr.Write([]byte(response))
	})

	return mu
}
