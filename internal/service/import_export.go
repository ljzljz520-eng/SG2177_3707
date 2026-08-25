package service

import (
	"context"
	"fmt"
	"path/filepath"

	"equipmentlending/internal/audit"
	"equipmentlending/internal/equipment"
	"equipmentlending/internal/persistence"
)

type PageResult struct {
	Records []equipment.EquipmentRecord
	Offset  int
	Limit   int
	Total   int
	More    bool
}

func (s *Service) Export(ctx context.Context, path string) error {
	backup, err := s.store.ExportBackup(ctx)
	if err != nil {
		return err
	}
	if err := persistence.WriteBackup(path, backup); err != nil {
		return err
	}
	return s.store.SaveAudit(ctx, persistence.AuditEntry{ID: nextAuditID(backup.Audit), Action: audit.ActionExported, CreatedAt: "export"})
}

func (s *Service) Import(ctx context.Context, path string) error {
	backup, err := persistence.ReadBackup(path)
	if err != nil {
		return err
	}
	if err := equipment.ValidateCollection(backup.Records); err != nil {
		return fmt.Errorf("backup validation: %w", err)
	}
	if err := s.store.ImportBackup(ctx, backup); err != nil {
		return err
	}
	return s.Load(ctx)
}

func (s *Service) Page(ctx context.Context, borrower, status string, offset, limit int) (PageResult, error) {
	records, err := s.Query(ctx, borrower, status)
	if err != nil {
		return PageResult{}, err
	}
	page := equipment.Page(records, offset, limit)
	return PageResult{Records: page, Offset: offset, Limit: limit, Total: len(records), More: offset+len(page) < len(records)}, nil
}

func (s *Service) Audit(ctx context.Context) ([]persistence.AuditEntry, string, error) {
	entries, err := s.store.LoadAudit(ctx)
	if err != nil {
		return nil, "", err
	}
	return entries, audit.Summary(entries), nil
}

func (s *Service) Inventory(ctx context.Context) (equipment.InventorySummary, error) {
	records, err := s.Query(ctx, "", "")
	if err != nil {
		return equipment.InventorySummary{}, err
	}
	return equipment.Summarize(records), nil
}

func (s *Service) SaveSetting(ctx context.Context, key, value string) error {
	if err := s.store.SetSetting(ctx, key, value); err != nil {
		return err
	}
	return s.store.SaveAudit(ctx, persistence.AuditEntry{ID: nextAuditID(nil), Action: "setting_changed", CreatedAt: key})
}

func nextAuditID(entries []persistence.AuditEntry) int64 {
	var highest int64
	for _, entry := range entries {
		if entry.ID > highest {
			highest = entry.ID
		}
	}
	if highest == 0 {
		return 1
	}
	return highest + 1
}

func archivePath(base, name string) string {
	return filepath.Join(base, name)
}
