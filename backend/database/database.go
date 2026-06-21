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
	}
	for _, s := range statements {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}
