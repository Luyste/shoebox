package web

import (
	"database/sql"
	"embed"

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
