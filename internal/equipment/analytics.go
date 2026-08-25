package equipment

import (
	"errors"
	"sort"
	"strings"
)

type InventorySummary struct {
	Total       int
	Available   int
	Borrowed    int
	Returned    int
	Borrowers   int
	Earliest    string
	Latest      string
	NumberOrder bool
}

func Summarize(records []EquipmentRecord) InventorySummary {
	counts := CountByStatus(records)
	summary := InventorySummary{Total: len(records), Available: counts[StatusAvailable], Borrowed: counts[StatusBorrowed], Returned: counts[StatusReturned], Borrowers: len(UniqueBorrowers(records)), NumberOrder: IsSortedByNumber(records)}
	for _, record := range records {
		if summary.Earliest == "" || DateBefore(record.BorrowDate, summary.Earliest) {
			summary.Earliest = record.BorrowDate
		}
		if summary.Latest == "" || DateAfter(record.BorrowDate, summary.Latest) {
			summary.Latest = record.BorrowDate
		}
	}
	return summary
}

func FindByID(records []EquipmentRecord, id int64) (EquipmentRecord, bool) {
	for _, record := range records {
		if record.ID == id {
			return record.Clone(), true
		}
	}
	return EquipmentRecord{}, false
}

func FindByName(records []EquipmentRecord, name string) []EquipmentRecord {
	needle := strings.ToLower(strings.TrimSpace(name))
	result := make([]EquipmentRecord, 0)
	for _, record := range records {
		if strings.Contains(strings.ToLower(record.Name), needle) {
			result = append(result, record.Clone())
		}
	}
	return result
}

func Page(records []EquipmentRecord, offset, limit int) []EquipmentRecord {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || offset >= len(records) {
		return []EquipmentRecord{}
	}
	end := offset + limit
	if end > len(records) {
		end = len(records)
	}
	return cloneRecords(records[offset:end])
}

func MergeRecords(primary, secondary []EquipmentRecord) ([]EquipmentRecord, error) {
	merged := make([]EquipmentRecord, 0, len(primary)+len(secondary))
	seen := make(map[string]bool, len(primary)+len(secondary))
	for _, collection := range [][]EquipmentRecord{primary, secondary} {
		for _, record := range collection {
			if record.EquipmentNumber == "" {
				return nil, errors.New("cannot merge record without number")
			}
			if seen[record.EquipmentNumber] {
				return nil, errors.New("duplicate equipment number in merge")
			}
			seen[record.EquipmentNumber] = true
			merged = append(merged, record.Clone())
		}
	}
	return merged, nil
}

func SortBorrowers(records []EquipmentRecord) []EquipmentRecord {
	result := cloneRecords(records)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Borrower == result[j].Borrower {
			return CompareNumber(result[i], result[j]) < 0
		}
		return result[i].Borrower < result[j].Borrower
	})
	return result
}

func ReplaceStatus(records []EquipmentRecord, from, to Status) ([]EquipmentRecord, int, error) {
	if !ValidStatus(from) || !ValidStatus(to) {
		return nil, 0, errors.New("invalid replacement status")
	}
	result := cloneRecords(records)
	changed := 0
	for index := range result {
		if result[index].Status == from {
			result[index].Status = to
			changed++
		}
	}
	return result, changed, nil
}

func ValidateCollection(records []EquipmentRecord) error {
	seen := make(map[string]bool, len(records))
	for _, record := range records {
		if err := record.Validate(); err != nil {
			return err
		}
		if seen[record.EquipmentNumber] {
			return errors.New("duplicate equipment number")
		}
		seen[record.EquipmentNumber] = true
	}
	return nil
}
