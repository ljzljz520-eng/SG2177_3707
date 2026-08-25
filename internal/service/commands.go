package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"equipmentlending/internal/equipment"
)

type Command struct {
	Name   string
	Number string
	Value  string
}

func (s *Service) Execute(ctx context.Context, command Command) (string, error) {
	name := strings.ToLower(strings.TrimSpace(command.Name))
	switch name {
	case "create":
		record, err := s.Create(ctx, CreateRequest{EquipmentNumber: command.Number, Name: command.Value, Borrower: "console", BorrowDate: "2026-01-01"})
		if err != nil {
			return "", err
		}
		return record.Summary(), nil
	case "delete":
		record, err := s.Delete(ctx, command.Number)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("deleted %s", record.EquipmentNumber), nil
	case "status":
		if err := s.UpdateStatus(ctx, command.Number, command.Value); err != nil {
			return "", err
		}
		return "status updated", nil
	case "sort":
		records, err := s.Sort(ctx, command.Value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("sorted %d records", len(records)), nil
	case "count":
		records, err := s.Query(ctx, "", "")
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d records", len(records)), nil
	default:
		return "", errors.New("unknown command")
	}
}

func CommandNames() []string {
	return []string{"create", "delete", "status", "sort", "count"}
}

func IsKnownCommand(name string) bool {
	for _, candidate := range CommandNames() {
		if candidate == strings.ToLower(strings.TrimSpace(name)) {
			return true
		}
	}
	return false
}

func CanDelete(record equipment.EquipmentRecord) bool {
	return record.Status != equipment.StatusBorrowed
}
