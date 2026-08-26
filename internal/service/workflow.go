package service

import (
	"context"
	"fmt"

	"equipmentlending/internal/equipment"
	"equipmentlending/internal/persistence"
)

type WorkflowResult struct {
	Created equipment.EquipmentRecord
	Records []equipment.EquipmentRecord
	Report  string
}

func ExecuteBorrowWorkflow(ctx context.Context, store *persistence.EquipmentStore, requests []CreateRequest) (WorkflowResult, error) {
	service := New(store)
	if err := service.Load(ctx); err != nil {
		return WorkflowResult{}, err
	}
	result := WorkflowResult{}
	for index, request := range requests {
		record, err := service.Create(ctx, request)
		if err != nil {
			return WorkflowResult{}, fmt.Errorf("create request %d: %w", index, err)
		}
		result.Created = record
	}
	result.Records = service.Snapshot()
	return result, nil
}

func ExecuteQueryWorkflow(ctx context.Context, store *persistence.EquipmentStore, borrower string) (WorkflowResult, error) {
	service := New(store)
	if err := service.Load(ctx); err != nil {
		return WorkflowResult{}, err
	}
	records, err := service.Query(ctx, borrower, "")
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Records: records}, nil
}

func ExecuteSortWorkflow(ctx context.Context, store *persistence.EquipmentStore, mode string) (WorkflowResult, error) {
	service := New(store)
	if err := service.Load(ctx); err != nil {
		return WorkflowResult{}, err
	}
	records, err := service.Sort(ctx, mode)
	if err != nil {
		return WorkflowResult{}, err
	}
	return WorkflowResult{Records: records}, nil
}

func EnsureComplete(records []equipment.EquipmentRecord, expected int) error {
	if len(records) != expected {
		return fmt.Errorf("expected %d records, got %d", expected, len(records))
	}
	seen := make(map[string]bool, len(records))
	for _, record := range records {
		if seen[record.EquipmentNumber] {
			return fmt.Errorf("duplicate equipment number %s", record.EquipmentNumber)
		}
		seen[record.EquipmentNumber] = true
	}
	return nil
}
