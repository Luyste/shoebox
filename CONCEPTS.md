# Concepts left to cover

Running list of what the project still needs before it is a functioning website,
and which Go/web concepts each part exercises.

## Broken right now

- **Thumbnail route.** `internal/web/views/grid.templ` points `<img>` at
  `/thumbs/<id>.jpg`, but `internal/web/router.go` only registers `/static/` and
  `/`. Every thumbnail 404s. Concept: serving files from disk (`os.DirFS`,
  `http.FileServerFS`, `http.StripPrefix`) next to the embedded static assets.
- **Download handler.** The TODO in `grid.templ` is still unwired. Concept:
  streaming a file response, `Content-Disposition`, `http.ServeFile`/`ServeContent`.
- **Pagination.** `Home` reads `?page=`, but no template emits `hx-get`, so no
  page past 0 is reachable. Concept: htmx partial swaps and infinite scroll
  (`hx-trigger="revealed"`, `hx-swap="afterend"`).
- **Indexer is unreachable from the app.** `indexer.New` is only called from
  `internal/web/seed_scratch_test.go`. There is no CLI subcommand and no upload
  endpoint, so there is no user-facing way to get media into the library.

## Concepts not touched yet

- **Writes from the browser.** Everything so far is read-only. Needs form
  handling (`r.ParseForm`), multipart uploads (`r.ParseMultipartForm`,
  `r.FormFile`), size limits (`http.MaxBytesReader`), and writing to the
  library tmp dir before moving into originals.
- **Server lifecycle.** `main.go` calls `http.ListenAndServe` with no timeouts.
  Needs an explicit `http.Server` with `ReadTimeout`/`WriteTimeout`/`IdleTimeout`
  and graceful shutdown on SIGINT via `signal.NotifyContext` + `srv.Shutdown`.
- **Request context.** Handlers use `db.Query`, not `db.QueryContext`, so a
  cancelled request keeps the query running.
- **Handler tests.** No `httptest` coverage of the web layer; the only tests are
  for db, media, and indexer.
- **Error handling.** Handlers currently return raw `err.Error()` as a 500,
  which leaks internals. Needs a small error-rendering helper or middleware.

## Minor / later

- `home.go` orders by `created_at`, but the only index in `0001_init.sql` is on
  `taken_at DESC, id DESC`. The index goes unused. Irrelevant until the table is
  large.
- Offset pagination will drift once media is inserted while paging. Keyset
  pagination on `(taken_at, id)` fixes it, and the index already exists.

## Suggested order

1. Thumbnail route (unblocks the whole grid visually)
2. htmx pagination
3. `index` subcommand in `main.go`
4. Download handler
5. Upload endpoint, then server timeouts and graceful shutdown
