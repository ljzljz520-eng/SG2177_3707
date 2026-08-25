package equipment

import "testing"

func TestEquipmentListLifecycle(t *testing.T) {
	list := NewList()
	first := NewRecord("A-02", "Meter", "Lin", "2026-01-02", StatusAvailable)
	second := NewRecord("A-01", "Scope", "Mei", "2026-01-03", StatusBorrowed)
	list.Append(first)
	list.Append(second)
	if list.Len() != 2 || list.Head().Value.EquipmentNumber != "A-02" || list.Tail().Value.EquipmentNumber != "A-01" {
		t.Fatalf("unexpected list endpoints")
	}
	if err := list.UpdateStatus("A-02", StatusBorrowed); err != nil {
		t.Fatal(err)
	}
	removed, ok := list.RemoveByNumber("A-01")
	if !ok || removed.Name != "Scope" || list.Len() != 1 {
		t.Fatalf("remove did not update list")
	}
	if err := list.ValidateLinks(); err != nil {
		t.Fatal(err)
	}
}

func TestRecordValidation(t *testing.T) {
	valid := NewRecord("A-1", "Meter", "Lin", "2026-02-03", StatusAvailable)
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	invalid := NewRecord("", "Meter", "Lin", "2026-02-03", StatusAvailable)
	if err := invalid.Validate(); err == nil {
		t.Fatal("expected missing number error")
	}
}
