package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"github.com/luyste/shoebox/internal/db"
	"github.com/luyste/shoebox/internal/library"
	"github.com/luyste/shoebox/internal/web"
)

//go:embed static
var staticFS embed.FS

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "./data", "library root")
	flag.Parse()

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

	mux := web.New(conn, lib, &staticFS).Router()
	h := web.Logger(mux)

	log.Fatal(http.ListenAndServe(":3000", h))
}
