package equipment

import "strings"

func FilterByBorrower(records []EquipmentRecord, borrower string) []EquipmentRecord {
	needle := strings.ToLower(strings.TrimSpace(borrower))
	filtered := make([]EquipmentRecord, 0, len(records))
	for _, record := range records {
		if needle == "" || strings.Contains(strings.ToLower(record.Borrower), needle) {
			filtered = append(filtered, record.Clone())
		}
	}
	return filtered
}

func FilterByStatus(records []EquipmentRecord, status Status) []EquipmentRecord {
	filtered := make([]EquipmentRecord, 0, len(records))
	for _, record := range records {
		if record.Status == status {
			filtered = append(filtered, record.Clone())
		}
	}
	return filtered
}

func CountByStatus(records []EquipmentRecord) map[Status]int {
	counts := map[Status]int{StatusAvailable: 0, StatusBorrowed: 0, StatusReturned: 0}
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}

func UniqueBorrowers(records []EquipmentRecord) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)
	for _, record := range records {
		if _, exists := seen[record.Borrower]; exists {
			continue
		}
		seen[record.Borrower] = struct{}{}
		result = append(result, record.Borrower)
	}
	return result
}
