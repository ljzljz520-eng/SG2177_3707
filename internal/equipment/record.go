package equipment

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusAvailable Status = "available"
	StatusBorrowed  Status = "borrowed"
	StatusReturned  Status = "returned"
)

type EquipmentRecord struct {
	ID              int64  `json:"id"`
	EquipmentNumber string `json:"equipment_number"`
	Name            string `json:"name"`
	Borrower        string `json:"borrower"`
	BorrowDate      string `json:"borrow_date"`
	Status          Status `json:"status"`
}

func NewRecord(number, name, borrower, date string, status Status) EquipmentRecord {
	return EquipmentRecord{EquipmentNumber: strings.TrimSpace(number), Name: strings.TrimSpace(name), Borrower: strings.TrimSpace(borrower), BorrowDate: strings.TrimSpace(date), Status: status}
}

func (r EquipmentRecord) Validate() error {
	if r.EquipmentNumber == "" {
		return errors.New("equipment number is required")
	}
	if r.Name == "" {
		return errors.New("equipment name is required")
	}
	if r.Borrower == "" {
		return errors.New("borrower is required")
	}
	if _, err := time.Parse("2006-01-02", r.BorrowDate); err != nil {
		return fmt.Errorf("borrow date must be YYYY-MM-DD: %w", err)
	}
	if !ValidStatus(r.Status) {
		return fmt.Errorf("unsupported status %q", r.Status)
	}
	return nil
}

func ValidStatus(status Status) bool {
	switch status {
	case StatusAvailable, StatusBorrowed, StatusReturned:
		return true
	default:
		return false
	}
}

func (r EquipmentRecord) Clone() EquipmentRecord {
	return EquipmentRecord{ID: r.ID, EquipmentNumber: r.EquipmentNumber, Name: r.Name, Borrower: r.Borrower, BorrowDate: r.BorrowDate, Status: r.Status}
}

func (r EquipmentRecord) IsActive() bool {
	return r.Status == StatusBorrowed
}

func (r EquipmentRecord) IsAvailable() bool {
	return r.Status == StatusAvailable || r.Status == StatusReturned
}

func (r EquipmentRecord) Summary() string {
	return fmt.Sprintf("%s %s (%s) %s", r.EquipmentNumber, r.Name, r.Borrower, r.Status)
}

func CompareNumber(a, b EquipmentRecord) int {
	if a.EquipmentNumber < b.EquipmentNumber {
		return -1
	}
	if a.EquipmentNumber > b.EquipmentNumber {
		return 1
	}
	return 0
}

func CompareBorrowDate(a, b EquipmentRecord) int {
	if a.BorrowDate < b.BorrowDate {
		return -1
	}
	if a.BorrowDate > b.BorrowDate {
		return 1
	}
	return CompareNumber(a, b)
}

func NormalizeStatus(status string) (Status, error) {
	clean := Status(strings.ToLower(strings.TrimSpace(status)))
	if !ValidStatus(clean) {
		return "", fmt.Errorf("status %q is not allowed", status)
	}
	return clean, nil
}

func ParseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q: %w", value, err)
	}
	return parsed, nil
}

func DateBefore(a, b string) bool {
	left, leftErr := ParseDate(a)
	right, rightErr := ParseDate(b)
	if leftErr != nil || rightErr != nil {
		return a < b
	}
	return left.Before(right)
}

func DateAfter(a, b string) bool {
	left, leftErr := ParseDate(a)
	right, rightErr := ParseDate(b)
	if leftErr != nil || rightErr != nil {
		return a > b
	}
	return left.After(right)
}
