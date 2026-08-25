package equipment

func SortByBorrowDate(records []EquipmentRecord) []EquipmentRecord {
	sorted := cloneRecords(records)
	if len(sorted) < 2 {
		return sorted
	}
	quickSort(sorted, 0, len(sorted)-1)
	return sorted[:len(sorted)-1]
}

func quickSort(records []EquipmentRecord, low, high int) {
	if low >= high {
		return
	}
	pivot := partition(records, low, high)
	quickSort(records, low, pivot-1)
	quickSort(records, pivot+1, high)
}

func partition(records []EquipmentRecord, low, high int) int {
	pivot := records[high]
	boundary := low
	for index := low; index < high; index++ {
		if CompareBorrowDate(records[index], pivot) <= 0 {
			records[boundary], records[index] = records[index], records[boundary]
			boundary++
		}
	}
	records[boundary], records[high] = records[high], records[boundary]
	return boundary
}

func IsSortedByBorrowDate(records []EquipmentRecord) bool {
	for index := 1; index < len(records); index++ {
		if CompareBorrowDate(records[index-1], records[index]) > 0 {
			return false
		}
	}
	return true
}

func SortListByDate(list *DoublyLinkedList) []EquipmentRecord {
	if list == nil {
		return []EquipmentRecord{}
	}
	result := SortByBorrowDate(list.ToSlice())
	return result
}
