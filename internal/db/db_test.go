package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpen_CreatesDatabase(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "profiles.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open(%q) error: %v", dbPath, err)
	}
	defer store.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("database file not created at %s", dbPath)
	}
}

func TestOpen_CreatesTablesOnFirstRun(t *testing.T) {
	dir := t.TempDir()
	store, err := Open(filepath.Join(dir, "profiles.db"))
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer store.Close()

	tables := []string{"profiles", "identifiers", "settings_global", "settings_user", "saved_locations", "usage_stats"}
	for _, table := range tables {
		var name string
		err := store.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestOpen_WALMode(t *testing.T) {
	dir := t.TempDir()
	store, err := Open(filepath.Join(dir, "profiles.db"))
	if err != nil {
		t.Fatalf("Open error: %v", err)
	}
	defer store.Close()

	var mode string
	err = store.db.QueryRow("PRAGMA journal_mode").Scan(&mode)
	if err != nil {
		t.Fatalf("PRAGMA journal_mode error: %v", err)
	}
	if mode != "wal" {
		t.Errorf("journal_mode = %q, want %q", mode, "wal")
	}
}

func TestOpen_IdempotentMigration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profiles.db")

	store1, err := Open(path)
	if err != nil {
		t.Fatalf("first Open error: %v", err)
	}
	store1.Close()

	store2, err := Open(path)
	if err != nil {
		t.Fatalf("second Open error: %v", err)
	}
	defer store2.Close()
}
