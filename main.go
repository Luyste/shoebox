package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luyste/shoebox/internal/db"
	"github.com/luyste/shoebox/internal/library"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "./data", "library root")

	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	lib, err := library.New(dir)
	if err != nil {
		log.Fatalf("open library failed: %v", err)
	}

	conn, err := db.Connect(lib.DbPath())
	if err != nil {
		log.Fatalf("database connection establish failed: %v", err)
	}

	err = db.Migrate(conn)
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	http.ListenAndServe(":3000", r)
}
