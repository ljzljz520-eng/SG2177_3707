package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"equipmentlending/internal/equipment"
	_ "modernc.org/sqlite"
)

type AuditEntry struct {
	ID        int64
	Action    string
	RecordID  int64
	CreatedAt string
}

type EquipmentStore struct {
	path string
	db   *sql.DB
}

func Open(path string) (*EquipmentStore, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	store := &EquipmentStore{path: path, db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := createSchema(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}
	return store, nil
}

func (s *EquipmentStore) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

func (s *EquipmentStore) DB() *sql.DB {
	if s == nil {
		return nil
	}
	return s.db
}

func (s *EquipmentStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *EquipmentStore) Healthy(ctx context.Context) error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	var result int
	if err := s.db.QueryRowContext(ctx, `SELECT 1`).Scan(&result); err != nil {
		return err
	}
	if result != 1 {
		return errors.New("unexpected health result")
	}
	return nil
}

func (s *EquipmentStore) SaveRecords(ctx context.Context, records []equipment.EquipmentRecord) error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	rollback := func() error { return tx.Rollback() }
	if _, err := tx.ExecContext(ctx, `DELETE FROM equipment_records`); err != nil {
		return rollback()
	}
	for _, record := range records {
		if err := record.Validate(); err != nil {
			return rollback()
		}
		_, err := tx.ExecContext(ctx, `INSERT INTO equipment_records(id, equipment_number, name, borrower, borrow_date, status, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)`, record.ID, record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, record.Status, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return rollback()
		}
	}
	return tx.Commit()
}

func (s *EquipmentStore) LoadRecords(ctx context.Context) ([]equipment.EquipmentRecord, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("store is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, equipment_number, name, borrower, borrow_date, status FROM equipment_records ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]equipment.EquipmentRecord, 0)
	for rows.Next() {
		var record equipment.EquipmentRecord
		var status string
		if err := rows.Scan(&record.ID, &record.EquipmentNumber, &record.Name, &record.Borrower, &record.BorrowDate, &status); err != nil {
			return nil, err
		}
		record.Status = equipment.Status(status)
		result = append(result, record)
	}
	return result, rows.Err()
}

func (s *EquipmentStore) SaveAudit(ctx context.Context, entry AuditEntry) error {
	if entry.CreatedAt == "" {
		entry.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_entries(id, action, record_id, created_at) VALUES(?, ?, ?, ?)`, entry.ID, entry.Action, entry.RecordID, entry.CreatedAt)
	return err
}

func (s *EquipmentStore) LoadAudit(ctx context.Context) ([]AuditEntry, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, action, record_id, created_at FROM audit_entries ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	entries := make([]AuditEntry, 0)
	for rows.Next() {
		var entry AuditEntry
		if err := rows.Scan(&entry.ID, &entry.Action, &entry.RecordID, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (s *EquipmentStore) SetSetting(ctx context.Context, key, value string) error {
	if key == "" {
		return errors.New("setting key is required")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO settings(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}

func (s *EquipmentStore) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *EquipmentStore) Ready(ctx context.Context) (bool, error) {
	return schemaReady(ctx, s.db)
}

func (s *EquipmentStore) Tables(ctx context.Context) ([]string, error) {
	return tableNames(ctx, s.db)
}
