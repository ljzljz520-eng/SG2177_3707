package reporting

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"equipmentlending/internal/equipment"
)

type Report struct {
	Title   string
	Rows    []string
	Summary string
}

func Build(records []equipment.EquipmentRecord, title string) Report {
	ordered := equipment.SortByNumber(records)
	rows := make([]string, 0, len(ordered))
	for index, record := range ordered {
		rows = append(rows, fmt.Sprintf("%02d | %s | %s | %s | %s | %s", index+1, record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, record.Status))
	}
	counts := equipment.CountByStatus(records)
	summary := fmt.Sprintf("total=%d available=%d borrowed=%d returned=%d", len(records), counts[equipment.StatusAvailable], counts[equipment.StatusBorrowed], counts[equipment.StatusReturned])
	return Report{Title: title, Rows: rows, Summary: summary}
}

func (r Report) String() string {
	parts := []string{r.Title}
	parts = append(parts, r.Rows...)
	parts = append(parts, r.Summary)
	return strings.Join(parts, "\n")
}

func (r Report) Write(w io.Writer) error {
	if w == nil {
		return fmt.Errorf("writer is nil")
	}
	_, err := io.WriteString(w, r.String()+"\n")
	return err
}

func FormatRecords(records []equipment.EquipmentRecord) string {
	return Build(records, "Equipment Lending Records").String()
}

func GroupByBorrower(records []equipment.EquipmentRecord) map[string][]equipment.EquipmentRecord {
	groups := make(map[string][]equipment.EquipmentRecord)
	for _, record := range records {
		groups[record.Borrower] = append(groups[record.Borrower], record.Clone())
	}
	for borrower := range groups {
		sort.SliceStable(groups[borrower], func(i, j int) bool { return groups[borrower][i].BorrowDate < groups[borrower][j].BorrowDate })
	}
	return groups
}
