package equipment

import "testing"

func TestSortByNumber(t *testing.T) {
	records := []EquipmentRecord{
		NewRecord("C-3", "C", "A", "2026-01-03", StatusAvailable),
		NewRecord("A-1", "A", "B", "2026-01-01", StatusBorrowed),
		NewRecord("B-2", "B", "C", "2026-01-02", StatusReturned),
	}
	sorted := SortByNumber(records)
	if len(sorted) != 3 || !IsSortedByNumber(sorted) {
		t.Fatalf("number sort failed: %#v", sorted)
	}
}

func TestEquipmentSortKeepsAll(t *testing.T) {
	records := []EquipmentRecord{
		NewRecord("A-1", "A", "A", "2026-03-03", StatusBorrowed),
		NewRecord("A-2", "B", "B", "2026-01-01", StatusAvailable),
		NewRecord("A-3", "C", "C", "2026-02-02", StatusReturned),
	}
	sorted := SortByBorrowDate(records)
	if len(sorted) != len(records) {
		t.Fatalf("date sort lost records: got %d want %d", len(sorted), len(records))
	}
	if !IsSortedByBorrowDate(sorted) {
		t.Fatalf("date sort order is invalid")
	}
}
