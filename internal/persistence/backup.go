package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"equipmentlending/internal/equipment"
)

type Backup struct {
	Records  []equipment.EquipmentRecord `json:"records"`
	Audit    []AuditEntry                `json:"audit"`
	Settings map[string]string           `json:"settings"`
}

func (s *EquipmentStore) ExportBackup(ctx context.Context) (Backup, error) {
	records, err := s.LoadRecords(ctx)
	if err != nil {
		return Backup{}, err
	}
	audit, err := s.LoadAudit(ctx)
	if err != nil {
		return Backup{}, err
	}
	settings := make(map[string]string)
	for _, key := range []string{"sort_mode", "site"} {
		value, getErr := s.GetSetting(ctx, key)
		if getErr != nil {
			return Backup{}, getErr
		}
		if value != "" {
			settings[key] = value
		}
	}
	return Backup{Records: records, Audit: audit, Settings: settings}, nil
}

func WriteBackup(path string, backup Backup) error {
	if path == "" {
		return errors.New("backup path is required")
	}
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}

func ReadBackup(path string) (Backup, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Backup{}, err
	}
	var backup Backup
	if err := json.Unmarshal(data, &backup); err != nil {
		return Backup{}, err
	}
	return backup, nil
}

func (s *EquipmentStore) ImportBackup(ctx context.Context, backup Backup) error {
	if err := s.SaveRecords(ctx, backup.Records); err != nil {
		return err
	}
	for _, entry := range backup.Audit {
		if err := s.SaveAudit(ctx, entry); err != nil {
			return err
		}
	}
	for key, value := range backup.Settings {
		if err := s.SetSetting(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}
