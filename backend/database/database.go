package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Enable foreign keys for SQLite (needed for ON DELETE CASCADE)
	_, _ = db.Exec("PRAGMA foreign_keys = ON")

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS versions (
			version TEXT PRIMARY KEY,
			date TEXT NOT NULL,
			size TEXT NOT NULL DEFAULT '0 MB',
			notes TEXT NOT NULL DEFAULT '[]',
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS stats (
			date TEXT PRIMARY KEY,
			launches INTEGER NOT NULL DEFAULT 0,
			updates INTEGER NOT NULL DEFAULT 0,
			traffic INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS channels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			enabled INTEGER NOT NULL DEFAULT 1,
			weight INTEGER NOT NULL DEFAULT 0,
			config TEXT NOT NULL DEFAULT '{}',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS version_changes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version TEXT NOT NULL,
			action TEXT NOT NULL,
			path TEXT NOT NULL,
			FOREIGN KEY (version) REFERENCES versions(version) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS version_file_snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version TEXT NOT NULL,
			path TEXT NOT NULL,
			hash TEXT NOT NULL,
			FOREIGN KEY (version) REFERENCES versions(version) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshots_version ON version_file_snapshots(version)`,
	}
	for _, s := range statements {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}
