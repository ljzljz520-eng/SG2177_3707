package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"equipmentlending/internal/equipment"
)

func (s *EquipmentStore) CountRecords(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM equipment_records`).Scan(&count)
	return count, err
}

func (s *EquipmentStore) FindRecord(ctx context.Context, number string) (equipment.EquipmentRecord, error) {
	var record equipment.EquipmentRecord
	var status string
	err := s.db.QueryRowContext(ctx, `SELECT id, equipment_number, name, borrower, borrow_date, status FROM equipment_records WHERE equipment_number=?`, number).Scan(&record.ID, &record.EquipmentNumber, &record.Name, &record.Borrower, &record.BorrowDate, &status)
	if err == sql.ErrNoRows {
		return equipment.EquipmentRecord{}, errors.New("record not found")
	}
	if err != nil {
		return equipment.EquipmentRecord{}, err
	}
	record.Status = equipment.Status(status)
	return record, nil
}

func (s *EquipmentStore) DeleteRecord(ctx context.Context, number string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM equipment_records WHERE equipment_number=?`, number)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (s *EquipmentStore) UpdateRecordStatus(ctx context.Context, number string, status equipment.Status) error {
	if !equipment.ValidStatus(status) {
		return errors.New("invalid status")
	}
	result, err := s.db.ExecContext(ctx, `UPDATE equipment_records SET status=? WHERE equipment_number=?`, status, number)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("record not found")
	}
	return nil
}

func (s *EquipmentStore) VerifyIntegrity(ctx context.Context) error {
	var result string
	if err := s.db.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&result); err != nil {
		return err
	}
	if strings.ToLower(result) != "ok" {
		return fmt.Errorf("integrity check: %s", result)
	}
	return nil
}

func (s *EquipmentStore) Vacuum(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `VACUUM`); err != nil {
		return err
	}
	return nil
}

func (s *EquipmentStore) SetMetadata(ctx context.Context, key, value string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("metadata key is required")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO metadata(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}

func (s *EquipmentStore) GetMetadata(ctx context.Context, key string) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM metadata WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
