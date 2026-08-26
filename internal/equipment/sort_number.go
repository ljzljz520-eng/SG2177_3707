package equipment

func SortByNumber(records []EquipmentRecord) []EquipmentRecord {
	sorted := cloneRecords(records)
	for index := 1; index < len(sorted); index++ {
		current := sorted[index]
		position := index - 1
		for position >= 0 && CompareNumber(sorted[position], current) > 0 {
			sorted[position+1] = sorted[position]
			position--
		}
		sorted[position+1] = current
	}
	return sorted
}

func IsSortedByNumber(records []EquipmentRecord) bool {
	for index := 1; index < len(records); index++ {
		if CompareNumber(records[index-1], records[index]) > 0 {
			return false
		}
	}
	return true
}

func InsertSortedByNumber(list *DoublyLinkedList, record EquipmentRecord) {
	if list == nil || list.Len() == 0 {
		if list != nil {
			list.Append(record)
		}
		return
	}
	for node := list.Head(); node != nil; node = node.Next {
		if CompareNumber(record, node.Value) < 0 {
			previous := node.Prev
			newNode := &Node{Value: record.Clone(), Prev: previous, Next: node}
			if previous == nil {
				list.Prepend(record)
			} else {
				previous.Next = newNode
				node.Prev = newNode
				list.length++
			}
			return
		}
	}
	list.Append(record)
}

func cloneRecords(records []EquipmentRecord) []EquipmentRecord {
	clones := make([]EquipmentRecord, len(records))
	for index, record := range records {
		clones[index] = record.Clone()
	}
	return clones
}
