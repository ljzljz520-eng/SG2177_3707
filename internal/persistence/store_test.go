package persistence

import (
	"context"
	"path/filepath"
	"testing"

	"equipmentlending/internal/equipment"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "equipment.db")
	ctx := context.Background()
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	records := []equipment.EquipmentRecord{{ID: 1, EquipmentNumber: "P-1", Name: "Probe", Borrower: "Kai", BorrowDate: "2026-04-01", Status: equipment.StatusBorrowed}}
	if err := store.SaveRecords(ctx, records); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveAudit(ctx, AuditEntry{ID: 11, Action: "created", RecordID: 1, CreatedAt: "2026-04-01T00:00:00Z"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SetSetting(ctx, "sort_mode", "date"); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	loaded, err := reopened.LoadRecords(ctx)
	if err != nil || len(loaded) != 1 || loaded[0].EquipmentNumber != "P-1" {
		t.Fatalf("records did not survive reopen: %#v %v", loaded, err)
	}
	ready, err := reopened.Ready(ctx)
	if err != nil || !ready {
		t.Fatalf("schema is not ready: %v", err)
	}
}

func TestBackupRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "equipment.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	backup, err := store.ExportBackup(context.Background())
	if err != nil || backup.Records == nil {
		t.Fatalf("backup failed: %#v %v", backup, err)
	}
}
