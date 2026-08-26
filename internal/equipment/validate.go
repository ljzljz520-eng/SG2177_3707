package equipment

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var numberPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

type ValidationIssue struct {
	Field   string
	Message string
}

func ValidateRecordDetailed(record EquipmentRecord) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	if strings.TrimSpace(record.EquipmentNumber) == "" {
		issues = append(issues, ValidationIssue{Field: "equipment_number", Message: "required"})
	} else if !numberPattern.MatchString(record.EquipmentNumber) {
		issues = append(issues, ValidationIssue{Field: "equipment_number", Message: "contains unsupported characters"})
	}
	if strings.TrimSpace(record.Name) == "" {
		issues = append(issues, ValidationIssue{Field: "name", Message: "required"})
	}
	if strings.TrimSpace(record.Borrower) == "" {
		issues = append(issues, ValidationIssue{Field: "borrower", Message: "required"})
	}
	if strings.TrimSpace(record.BorrowDate) == "" {
		issues = append(issues, ValidationIssue{Field: "borrow_date", Message: "required"})
	} else if _, err := ParseDate(record.BorrowDate); err != nil {
		issues = append(issues, ValidationIssue{Field: "borrow_date", Message: "must use YYYY-MM-DD"})
	}
	if !ValidStatus(record.Status) {
		issues = append(issues, ValidationIssue{Field: "status", Message: "unsupported value"})
	}
	return issues
}

func ValidationError(record EquipmentRecord) error {
	issues := ValidateRecordDetailed(record)
	if len(issues) == 0 {
		return nil
	}
	parts := make([]string, 0, len(issues))
	for _, issue := range issues {
		parts = append(parts, fmt.Sprintf("%s %s", issue.Field, issue.Message))
	}
	return errors.New(strings.Join(parts, "; "))
}

func CleanRecord(record EquipmentRecord) EquipmentRecord {
	record.EquipmentNumber = strings.TrimSpace(record.EquipmentNumber)
	record.Name = strings.TrimSpace(record.Name)
	record.Borrower = strings.TrimSpace(record.Borrower)
	record.BorrowDate = strings.TrimSpace(record.BorrowDate)
	record.Status = Status(strings.ToLower(strings.TrimSpace(string(record.Status))))
	return record
}

func NormalizeRecords(records []EquipmentRecord) ([]EquipmentRecord, error) {
	result := make([]EquipmentRecord, 0, len(records))
	seen := make(map[string]bool, len(records))
	for _, source := range records {
		record := CleanRecord(source)
		if err := ValidationError(record); err != nil {
			return nil, err
		}
		if seen[record.EquipmentNumber] {
			return nil, fmt.Errorf("duplicate number %s", record.EquipmentNumber)
		}
		seen[record.EquipmentNumber] = true
		result = append(result, record)
	}
	return result, nil
}

func MissingFields(record EquipmentRecord) []string {
	missing := make([]string, 0)
	if strings.TrimSpace(record.EquipmentNumber) == "" {
		missing = append(missing, "equipment_number")
	}
	if strings.TrimSpace(record.Name) == "" {
		missing = append(missing, "name")
	}
	if strings.TrimSpace(record.Borrower) == "" {
		missing = append(missing, "borrower")
	}
	if strings.TrimSpace(record.BorrowDate) == "" {
		missing = append(missing, "borrow_date")
	}
	return missing
}

func CanTransition(from, to Status) bool {
	if !ValidStatus(from) || !ValidStatus(to) {
		return false
	}
	if from == to {
		return true
	}
	switch from {
	case StatusAvailable:
		return to == StatusBorrowed
	case StatusBorrowed:
		return to == StatusReturned
	case StatusReturned:
		return to == StatusAvailable || to == StatusBorrowed
	default:
		return false
	}
}

func Transition(record EquipmentRecord, to Status) (EquipmentRecord, error) {
	if !CanTransition(record.Status, to) {
		return EquipmentRecord{}, fmt.Errorf("cannot transition %s from %s to %s", record.EquipmentNumber, record.Status, to)
	}
	record.Status = to
	return record, nil
}
