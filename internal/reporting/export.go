package reporting

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"equipmentlending/internal/equipment"
)

func WriteCSV(w io.Writer, records []equipment.EquipmentRecord) error {
	if w == nil {
		return fmt.Errorf("writer is nil")
	}
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"id", "equipment_number", "name", "borrower", "borrow_date", "status"}); err != nil {
		return err
	}
	for _, record := range records {
		if err := writer.Write([]string{fmt.Sprint(record.ID), record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, string(record.Status)}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func Markdown(records []equipment.EquipmentRecord) string {
	lines := []string{"| Number | Name | Borrower | Borrow date | Status |", "|---|---|---|---|---|"}
	for _, record := range equipment.SortByNumber(records) {
		lines = append(lines, fmt.Sprintf("| %s | %s | %s | %s | %s |", record.EquipmentNumber, record.Name, record.Borrower, record.BorrowDate, record.Status))
	}
	return strings.Join(lines, "\n")
}

func ReportWithComplexity(records []equipment.EquipmentRecord) string {
	return FormatRecords(records) + "\n\nComplexity\n" + ComplexityText()
}

func StatusLegend() string {
	return "available=ready for lending; borrowed=currently out; returned=awaiting inspection"
}
