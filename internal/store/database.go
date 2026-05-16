package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Database struct {
	DB *sql.DB
}

func New(dataDir string) (*Database, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filepath.Join(dataDir, "webhix.db"))
	if err != nil {
		return nil, err
	}

	return &Database{
		DB: db,
	}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) Migrate(migrationsDir string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	return goose.Up(d.DB, migrationsDir)
}
