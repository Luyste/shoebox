package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrate(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	conn, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	if err := Migrate(conn); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
}

func TestDoubleMigrate(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	conn, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	if err := Migrate(conn); err != nil {
		t.Fatalf("first migrate failed: %v", err)
	}

	if err := Migrate(conn); err != nil {
		t.Fatalf("second migrate failed: %v", err)
	}
}

func TestChecksumDrift(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	conn, err := Connect(dbPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	if err := Migrate(conn); err != nil {
		t.Fatalf("frist migration failed: %v", err)
	}
	res, err := conn.Exec("UPDATE schema_migrations SET checksum = 'corrupt' WHERE file_name = '0001_init.sql'")
	if err != nil {
		t.Fatalf("updating migration checksum failed: %v", err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		t.Errorf("no rows affected by update: %v", affected)
	}

	err = Migrate(conn)
	if err == nil {
		t.Fatalf("migration after checksum invalidation was succesful")
	}
	if err != nil {
		if !strings.Contains(err.Error(), "0001_init.sql") {
			t.Errorf("expected error to contain filename, got: %v", err)
		}
	}

}
