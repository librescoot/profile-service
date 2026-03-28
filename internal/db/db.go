package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db   *sql.DB
	path string
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	pragmas := []string{
		"PRAGMA journal_mode=wal",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("exec %q: %w", p, err)
		}
	}

	store := &Store{db: db, path: path}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			cloud_user_id TEXT,
			avatar_type TEXT NOT NULL DEFAULT 'parametric',
			avatar_seed INTEGER,
			avatar_path TEXT,
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS identifiers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id TEXT NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
			type TEXT NOT NULL,
			value TEXT NOT NULL,
			created_at INTEGER NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_identifiers_type_value ON identifiers(type, value);

		CREATE TABLE IF NOT EXISTS settings_global (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS settings_user (
			profile_id TEXT NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
			key TEXT NOT NULL,
			value TEXT NOT NULL,
			PRIMARY KEY (profile_id, key)
		);

		CREATE TABLE IF NOT EXISTS saved_locations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			profile_id TEXT NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			type TEXT NOT NULL DEFAULT 'favorite',
			created_at INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		);

		CREATE TABLE IF NOT EXISTS usage_stats (
			profile_id TEXT PRIMARY KEY REFERENCES profiles(id) ON DELETE CASCADE,
			total_distance_m INTEGER NOT NULL DEFAULT 0,
			total_trips INTEGER NOT NULL DEFAULT 0,
			total_duration_s INTEGER NOT NULL DEFAULT 0,
			max_speed REAL NOT NULL DEFAULT 0,
			last_ride_at INTEGER
		);
	`)
	return err
}
