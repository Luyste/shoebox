package api

import (
	"database/sql"

	"github.com/luyste/shoebox/internal/library"
)

type API struct {
	db  *sql.DB
	lib *library.Library
}

func New(db *sql.DB, lib *library.Library) *API {
	return &API{db: db, lib: lib}
}
