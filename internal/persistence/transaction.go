package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"equipmentlending/internal/equipment"
)

type TransactionResult struct {
	RecordsWritten int
	AuditWritten   int
	Committed      bool
}

func (s *EquipmentStore) ReplaceAll(ctx context.Context, records []equipment.EquipmentRecord, entries []AuditEntry) (TransactionResult, error) {
	if s == nil || s.db == nil {
		return TransactionResult{}, errors.New("store is closed")
	}
	normalized, err := equipment.NormalizeRecords(records)
	if err != nil {
		return TransactionResult{}, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return TransactionResult{}, err
	}
	rollback := func(cause error) (TransactionResult, error) {
		_ = tx.Rollback()
		return TransactionResult{}, cause
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM equipment_records`); err != nil {
		return rollback(err)
	}
	for _, record := range normalized {
		if record.ID == 0 {
			_, err = tx.ExecContext(ctx, `INSERT INTO equipment_records(equipment_number, name, borrower, borrow_date, status, updated_at) VALUES(?, ?, ?, ?, ?, ?)`, record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, record.Status, time.Now().UTC().Format(time.RFC3339))
		} else {
			_, err = tx.ExecContext(ctx, `INSERT INTO equipment_records(id, equipment_number, name, borrower, borrow_date, status, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)`, record.ID, record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, record.Status, time.Now().UTC().Format(time.RFC3339))
		}
		if err != nil {
			return rollback(err)
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM audit_entries`); err != nil {
		return rollback(err)
	}
	for _, entry := range entries {
		if entry.CreatedAt == "" {
			entry.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO audit_entries(id, action, record_id, created_at) VALUES(?, ?, ?, ?)`, entry.ID, entry.Action, entry.RecordID, entry.CreatedAt); err != nil {
			return rollback(err)
		}
	}
	if err := tx.Commit(); err != nil {
		return TransactionResult{}, err
	}
	return TransactionResult{RecordsWritten: len(normalized), AuditWritten: len(entries), Committed: true}, nil
}

func (s *EquipmentStore) InTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if fn == nil {
		_ = tx.Rollback()
		return errors.New("transaction callback is required")
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *EquipmentStore) RecordAuditBatch(ctx context.Context, entries []AuditEntry) error {
	return s.InTransaction(ctx, func(tx *sql.Tx) error {
		for _, entry := range entries {
			if entry.Action == "" {
				return errors.New("audit action is required")
			}
			if entry.CreatedAt == "" {
				entry.CreatedAt = time.Now().UTC().Format(time.RFC3339)
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO audit_entries(id, action, record_id, created_at) VALUES(?, ?, ?, ?)`, entry.ID, entry.Action, entry.RecordID, entry.CreatedAt); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *EquipmentStore) DatabaseInfo(ctx context.Context) (map[string]string, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("store is closed")
	}
	info := make(map[string]string)
	var version string
	if err := s.db.QueryRowContext(ctx, `SELECT sqlite_version()`).Scan(&version); err != nil {
		return nil, err
	}
	info["sqlite_version"] = version
	count, err := s.CountRecords(ctx)
	if err != nil {
		return nil, err
	}
	info["record_count"] = fmt.Sprint(count)
	ready, err := s.Ready(ctx)
	if err != nil {
		return nil, err
	}
	if ready {
		info["schema"] = "ready"
	} else {
		info["schema"] = "not_ready"
	}
	return info, nil
}
