package db

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type MigrationFile struct {
	filename string
	bytes    []byte
	filePath string
}

func Migrate(db *sql.DB) error {
	err := ensureSchemaMigrationsTable(db)
	if err != nil {
		return fmt.Errorf("migrations table init failed: %w", err)
	}

	migrationFiles, err := readAndSortMigrations(&migrationsFS)
	if err != nil {
		return fmt.Errorf("migration file parsing failed: %w", err)
	}

	for _, mf := range migrationFiles {
		var diskChecksum string

		checksum := buildChecksum(mf.bytes)

		err = db.QueryRow(`SELECT checksum FROM schema_migrations WHERE file_name = ?`, mf.filename).Scan(&diskChecksum)

		switch {
		case err == nil:
			if diskChecksum != checksum {
				return fmt.Errorf("checksum of %s corrupted", mf.filePath)
			}
		case errors.Is(err, sql.ErrNoRows):
			err = applyMigration(db, string(mf.bytes), mf.filename, checksum)
			if err != nil {
				return fmt.Errorf("applying migration of %s failed: %w", mf.filePath, err)
			}
		default:
			return fmt.Errorf("something went wrong during migration of %s: %w", mf.filePath, err)
		}
	}
	return nil
}

func buildChecksum(content []byte) string {
	sha := sha256.Sum256(content)
	return hex.EncodeToString(sha[:])
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		file_name TEXT PRIMARY KEY,
		checksum TEXT NOT NULL,
		applied_at TEXT NOT NULL DEFAULT (datetime('now')))`); err != nil {
		return err
	}
	return nil
}

func readAndSortMigrations(embedfs *embed.FS) ([]MigrationFile, error) {
	migrationFiles := []MigrationFile{}

	files, err := fs.Glob(embedfs, "migrations/*.sql")
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	for _, file := range files {
		filename := strings.Split(file, "/")[1]
		bytes, err := fs.ReadFile(migrationsFS, file)
		if err != nil {
			return nil, fmt.Errorf("reading migration file %s: %w", file, err)
		}

		migrationFiles = append(migrationFiles, MigrationFile{
			filename: filename,
			bytes:    bytes,
			filePath: file,
		})
	}

	return migrationFiles, nil
}

func applyMigration(db *sql.DB, mgr string, filename string, checksum string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec(mgr)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO schema_migrations (file_name, checksum) VALUES (?, ?)`, filename, checksum)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
