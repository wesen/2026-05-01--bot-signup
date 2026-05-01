package database

import (
	"context"
	"path/filepath"
	"testing"
)

func openTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestOpenRunsMigrations(t *testing.T) {
	db := openTestDB(t)
	var count int
	if err := db.db.QueryRowContext(context.Background(), `SELECT count(*) FROM schema_migrations WHERE version = '001_initial.sql'`).Scan(&count); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected migration to be recorded once, got %d", count)
	}
}
