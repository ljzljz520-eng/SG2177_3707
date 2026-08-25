package persistence

import (
	"context"
	"database/sql"
)

const schemaVersion = 1

func createSchema(ctx context.Context, db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS equipment_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			equipment_number TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			borrower TEXT NOT NULL,
			borrow_date TEXT NOT NULL,
			status TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS audit_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT NOT NULL,
			record_id INTEGER NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS metadata (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		)`,
	}
	for _, statement := range statements {
		if _, err := db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	_, err := db.ExecContext(ctx, `INSERT INTO metadata(key, value) VALUES('schema_version', '1') ON CONFLICT(key) DO UPDATE SET value=excluded.value`)
	return err
}

func schemaReady(ctx context.Context, db *sql.DB) (bool, error) {
	var value int
	err := db.QueryRowContext(ctx, `SELECT CAST(value AS INTEGER) FROM metadata WHERE key='schema_version'`).Scan(&value)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return value == schemaVersion, nil
}

func tableNames(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		result = append(result, name)
	}
	return result, rows.Err()
}
