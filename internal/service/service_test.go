package service

import (
	"context"
	"path/filepath"
	"testing"

	"equipmentlending/internal/persistence"
)

func TestPrimaryWorkflow(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "primary.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	business := New(store)
	if err := business.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	record, err := business.Create(context.Background(), CreateRequest{EquipmentNumber: "W-1", Name: "Scale", Borrower: "Rui", BorrowDate: "2026-05-01", Status: "borrowed"})
	if err != nil || record.ID == 0 {
		t.Fatalf("create failed: %#v %v", record, err)
	}
	if err := business.UpdateStatus(context.Background(), "W-1", "returned"); err != nil {
		t.Fatal(err)
	}
	if _, err := business.Delete(context.Background(), "W-1"); err != nil {
		t.Fatal(err)
	}
}

func TestSecondaryWorkflow(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	business := New(store)
	if err := business.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, request := range []CreateRequest{{EquipmentNumber: "Q-1", Name: "A", Borrower: "Ming", BorrowDate: "2026-01-01"}, {EquipmentNumber: "Q-2", Name: "B", Borrower: "Ning", BorrowDate: "2026-01-02"}} {
		if _, err := business.Create(context.Background(), request); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := business.Query(context.Background(), "ming", "")
	if err != nil || len(rows) != 1 || rows[0].EquipmentNumber != "Q-1" {
		t.Fatalf("query failed: %#v %v", rows, err)
	}
}
