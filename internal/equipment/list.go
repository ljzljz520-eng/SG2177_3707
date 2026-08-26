package equipment

import "errors"

type Node struct {
	Value EquipmentRecord
	Prev  *Node
	Next  *Node
}

type DoublyLinkedList struct {
	head   *Node
	tail   *Node
	length int
}

func NewList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func (l *DoublyLinkedList) Len() int {
	if l == nil {
		return 0
	}
	return l.length
}

func (l *DoublyLinkedList) Head() *Node {
	if l == nil {
		return nil
	}
	return l.head
}

func (l *DoublyLinkedList) Tail() *Node {
	if l == nil {
		return nil
	}
	return l.tail
}

func (l *DoublyLinkedList) Append(record EquipmentRecord) *Node {
	if l == nil {
		return nil
	}
	node := &Node{Value: record.Clone(), Prev: l.tail}
	if l.tail == nil {
		l.head = node
	} else {
		l.tail.Next = node
	}
	l.tail = node
	l.length++
	return node
}

func (l *DoublyLinkedList) Prepend(record EquipmentRecord) *Node {
	if l == nil {
		return nil
	}
	node := &Node{Value: record.Clone(), Next: l.head}
	if l.head == nil {
		l.tail = node
	} else {
		l.head.Prev = node
	}
	l.head = node
	l.length++
	return node
}

func (l *DoublyLinkedList) RemoveByNumber(number string) (EquipmentRecord, bool) {
	node := l.FindByNumber(number)
	if node == nil {
		return EquipmentRecord{}, false
	}
	if node.Prev == nil {
		l.head = node.Next
	} else {
		node.Prev.Next = node.Next
	}
	if node.Next == nil {
		l.tail = node.Prev
	} else {
		node.Next.Prev = node.Prev
	}
	l.length--
	return node.Value.Clone(), true
}

func (l *DoublyLinkedList) FindByNumber(number string) *Node {
	if l == nil {
		return nil
	}
	for node := l.head; node != nil; node = node.Next {
		if node.Value.EquipmentNumber == number {
			return node
		}
	}
	return nil
}

func (l *DoublyLinkedList) UpdateStatus(number string, status Status) error {
	if !ValidStatus(status) {
		return errors.New("invalid status")
	}
	node := l.FindByNumber(number)
	if node == nil {
		return errors.New("record not found")
	}
	node.Value.Status = status
	return nil
}

func (l *DoublyLinkedList) ToSlice() []EquipmentRecord {
	if l == nil || l.length == 0 {
		return []EquipmentRecord{}
	}
	result := make([]EquipmentRecord, 0, l.length)
	for node := l.head; node != nil; node = node.Next {
		result = append(result, node.Value.Clone())
	}
	return result
}

func (l *DoublyLinkedList) ReverseSlice() []EquipmentRecord {
	if l == nil || l.length == 0 {
		return []EquipmentRecord{}
	}
	result := make([]EquipmentRecord, 0, l.length)
	for node := l.tail; node != nil; node = node.Prev {
		result = append(result, node.Value.Clone())
	}
	return result
}

func (l *DoublyLinkedList) Walk(fn func(EquipmentRecord) bool) int {
	if l == nil || fn == nil {
		return 0
	}
	count := 0
	for node := l.head; node != nil; node = node.Next {
		count++
		if !fn(node.Value.Clone()) {
			break
		}
	}
	return count
}

func (l *DoublyLinkedList) Clear() {
	if l == nil {
		return
	}
	l.head = nil
	l.tail = nil
	l.length = 0
}

func (l *DoublyLinkedList) Replace(records []EquipmentRecord) {
	l.Clear()
	for _, record := range records {
		l.Append(record)
	}
}

func ListFromRecords(records []EquipmentRecord) *DoublyLinkedList {
	list := NewList()
	list.Replace(records)
	return list
}

func (l *DoublyLinkedList) ValidateLinks() error {
	if l == nil {
		return errors.New("list is nil")
	}
	if l.length == 0 && (l.head != nil || l.tail != nil) {
		return errors.New("empty list has endpoints")
	}
	seen := 0
	var previous *Node
	for node := l.head; node != nil; node = node.Next {
		if node.Prev != previous {
			return errors.New("broken previous link")
		}
		previous = node
		seen++
		if seen > l.length {
			return errors.New("cycle detected")
		}
	}
	if previous != l.tail || seen != l.length {
		return errors.New("length or tail mismatch")
	}
	return nil
}
