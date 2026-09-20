package web

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/luyste/shoebox/internal/library"
	"github.com/luyste/shoebox/internal/web/views"
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

	mu.HandleFunc("GET /", w.Home)

	return mu
}

func (w *Web) mediaPage(page int) ([]views.Media, error) {
	const pageSize = 20
	offset := page * pageSize

	rows, err := w.db.Query("SELECT id, original_name FROM media ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []views.Media
	for rows.Next() {
		var m views.Media
		if err := rows.Scan(&m.Id, &m.OriginalName); err != nil {
			return nil, err
		}
		media = append(media, m)
	}

	return media, rows.Err()
}

func (w *Web) Home(wr http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 0
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	cards, err := w.mediaPage(page)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		views.Grid(cards).Render(r.Context(), wr)
	} else {
		views.Page(cards).Render(r.Context(), wr)
	}
}
