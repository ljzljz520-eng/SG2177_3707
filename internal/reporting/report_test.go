package reporting

import (
	"strings"
	"testing"

	"equipmentlending/internal/equipment"
)

func TestReportFormatting(t *testing.T) {
	records := []equipment.EquipmentRecord{equipment.NewRecord("R-2", "Tube", "Lin", "2026-01-02", equipment.StatusBorrowed), equipment.NewRecord("R-1", "Scale", "Mei", "2026-01-01", equipment.StatusAvailable)}
	text := FormatRecords(records)
	if !strings.Contains(text, "R-1") || !strings.Contains(text, "total=2") {
		t.Fatalf("report missing data: %s", text)
	}
	if !strings.Contains(ComplexityText(), "快速排序") {
		t.Fatal("complexity note missing")
	}
}
