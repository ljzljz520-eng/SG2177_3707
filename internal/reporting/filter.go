package reporting

import (
	"fmt"
	"sort"
	"strings"

	"equipmentlending/internal/equipment"
)

type QueryOptions struct {
	Borrower string
	Status   equipment.Status
	Contains string
	Sort     string
	Limit    int
}

func Apply(records []equipment.EquipmentRecord, options QueryOptions) []equipment.EquipmentRecord {
	result := make([]equipment.EquipmentRecord, 0, len(records))
	borrower := strings.ToLower(strings.TrimSpace(options.Borrower))
	contains := strings.ToLower(strings.TrimSpace(options.Contains))
	for _, record := range records {
		if borrower != "" && !strings.Contains(strings.ToLower(record.Borrower), borrower) {
			continue
		}
		if options.Status != "" && record.Status != options.Status {
			continue
		}
		if contains != "" && !strings.Contains(strings.ToLower(record.Name), contains) {
			continue
		}
		result = append(result, record.Clone())
	}
	switch options.Sort {
	case "date":
		result = equipment.SortByBorrowDate(result)
	case "borrower":
		result = equipment.SortBorrowers(result)
	default:
		result = equipment.SortByNumber(result)
	}
	if options.Limit > 0 && len(result) > options.Limit {
		result = result[:options.Limit]
	}
	return result
}

func Facets(records []equipment.EquipmentRecord) map[string][]string {
	result := map[string][]string{"borrowers": {}, "statuses": {}, "names": {}}
	borrowers := make(map[string]bool)
	statuses := make(map[string]bool)
	names := make(map[string]bool)
	for _, record := range records {
		borrowers[record.Borrower] = true
		statuses[string(record.Status)] = true
		names[record.Name] = true
	}
	for value := range borrowers {
		result["borrowers"] = append(result["borrowers"], value)
	}
	for value := range statuses {
		result["statuses"] = append(result["statuses"], value)
	}
	for value := range names {
		result["names"] = append(result["names"], value)
	}
	for key := range result {
		sort.Strings(result[key])
	}
	return result
}

func RenderTable(records []equipment.EquipmentRecord) string {
	if len(records) == 0 {
		return "No equipment records"
	}
	widths := []int{6, 8, 8, 8, 12, 10}
	for _, record := range records {
		values := []string{fmt.Sprint(record.ID), record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, string(record.Status)}
		for index, value := range values {
			if len(value) > widths[index] {
				widths[index] = len(value)
			}
		}
	}
	lines := []string{tableLine(widths), tableRow(widths, []string{"ID", "Number", "Name", "Borrower", "Borrow date", "Status"}), tableLine(widths)}
	for _, record := range records {
		lines = append(lines, tableRow(widths, []string{fmt.Sprint(record.ID), record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, string(record.Status)}))
	}
	lines = append(lines, tableLine(widths))
	return strings.Join(lines, "\n")
}

func tableLine(widths []int) string {
	parts := make([]string, len(widths))
	for index, width := range widths {
		parts[index] = strings.Repeat("-", width+2)
	}
	return "+" + strings.Join(parts, "+") + "+"
}

func tableRow(widths []int, values []string) string {
	parts := make([]string, len(widths))
	for index, width := range widths {
		value := values[index]
		parts[index] = " " + value + strings.Repeat(" ", width-len(value)+1)
	}
	return "|" + strings.Join(parts, "|") + "|"
}

func CompareReports(left, right Report) bool {
	if left.Title != right.Title || left.Summary != right.Summary || len(left.Rows) != len(right.Rows) {
		return false
	}
	for index := range left.Rows {
		if left.Rows[index] != right.Rows[index] {
			return false
		}
	}
	return true
}

func SafeTitle(value string) string {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return "Equipment Lending Records"
	}
	clean = strings.ReplaceAll(clean, "\n", " ")
	clean = strings.ReplaceAll(clean, "\r", " ")
	return clean
}

func BuildFiltered(records []equipment.EquipmentRecord, options QueryOptions, title string) Report {
	filtered := Apply(records, options)
	return Build(filtered, SafeTitle(title))
}

func StatusCountsLine(records []equipment.EquipmentRecord) string {
	counts := equipment.CountByStatus(records)
	return fmt.Sprintf("available:%d borrowed:%d returned:%d", counts[equipment.StatusAvailable], counts[equipment.StatusBorrowed], counts[equipment.StatusReturned])
}

func NumberIndex(records []equipment.EquipmentRecord) map[string]int {
	index := make(map[string]int, len(records))
	for position, record := range records {
		index[record.EquipmentNumber] = position
	}
	return index
}

func DistinctStatuses(records []equipment.EquipmentRecord) []equipment.Status {
	seen := make(map[equipment.Status]bool)
	result := make([]equipment.Status, 0)
	for _, record := range records {
		if !seen[record.Status] {
			seen[record.Status] = true
			result = append(result, record.Status)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}
