package service

import (
	"errors"
	"fmt"
	"strings"

	"equipmentlending/internal/equipment"
)

type CreateRequest struct {
	EquipmentNumber string
	Name            string
	Borrower        string
	BorrowDate      string
	Status          string
}

func (r CreateRequest) Record() (equipment.EquipmentRecord, error) {
	status := equipment.StatusAvailable
	if strings.TrimSpace(r.Status) != "" {
		parsed, err := equipment.NormalizeStatus(r.Status)
		if err != nil {
			return equipment.EquipmentRecord{}, err
		}
		status = parsed
	}
	record := equipment.NewRecord(r.EquipmentNumber, r.Name, r.Borrower, r.BorrowDate, status)
	if err := record.Validate(); err != nil {
		return equipment.EquipmentRecord{}, err
	}
	return record, nil
}

func ValidateNumber(number string) error {
	clean := strings.TrimSpace(number)
	if clean == "" {
		return errors.New("equipment number is required")
	}
	if len(clean) > 40 {
		return errors.New("equipment number is too long")
	}
	for _, char := range clean {
		if char == '/' || char == '\\' {
			return fmt.Errorf("equipment number contains forbidden character %q", char)
		}
	}
	return nil
}

func ValidateQuery(borrower, status string) error {
	if len(strings.TrimSpace(borrower)) > 100 {
		return errors.New("borrower query is too long")
	}
	if status != "" {
		if _, err := equipment.NormalizeStatus(status); err != nil {
			return err
		}
	}
	return nil
}

func NormalizeCreateRequest(request CreateRequest) CreateRequest {
	request.EquipmentNumber = strings.TrimSpace(request.EquipmentNumber)
	request.Name = strings.TrimSpace(request.Name)
	request.Borrower = strings.TrimSpace(request.Borrower)
	request.BorrowDate = strings.TrimSpace(request.BorrowDate)
	request.Status = strings.TrimSpace(strings.ToLower(request.Status))
	return request
}
