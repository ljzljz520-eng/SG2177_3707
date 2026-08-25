package main

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"equipmentlending/internal/persistence"
	"equipmentlending/internal/service"
)

func TestCommandHelpers(t *testing.T) {
	store, err := persistence.Open(filepath.Join(t.TempDir(), "cmd.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	business := service.New(store)
	if err := business.Load(context.Background()); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := runCreate(context.Background(), business, &output, service.CreateRequest{EquipmentNumber: "C-1", Name: "Caliper", Borrower: "Bo", BorrowDate: "2026-07-01"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "created") {
		t.Fatalf("unexpected output %s", output.String())
	}
}
