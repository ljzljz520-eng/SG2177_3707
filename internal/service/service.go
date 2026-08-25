package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"equipmentlending/internal/audit"
	"equipmentlending/internal/equipment"
	"equipmentlending/internal/persistence"
)

type Service struct {
	store *persistence.EquipmentStore
	list  *equipment.DoublyLinkedList
	mu    sync.RWMutex
	next  int64
}

func New(store *persistence.EquipmentStore) *Service {
	return &Service{store: store, list: equipment.NewList(), next: 1}
}

func (s *Service) Store() *persistence.EquipmentStore {
	return s.store
}

func (s *Service) Load(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	records, err := s.store.LoadRecords(ctx)
	if err != nil {
		return err
	}
	s.list.Replace(records)
	s.next = nextID(records)
	return nil
}

func (s *Service) Create(ctx context.Context, request CreateRequest) (equipment.EquipmentRecord, error) {
	request = NormalizeCreateRequest(request)
	if err := ValidateNumber(request.EquipmentNumber); err != nil {
		return equipment.EquipmentRecord{}, err
	}
	record, err := request.Record()
	if err != nil {
		return equipment.EquipmentRecord{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.list.FindByNumber(record.EquipmentNumber) != nil {
		return equipment.EquipmentRecord{}, errors.New("equipment number already exists")
	}
	record.ID = s.next
	s.next++
	s.list.Append(record)
	if err := s.persistLocked(ctx, audit.ActionCreated, record.ID); err != nil {
		return equipment.EquipmentRecord{}, err
	}
	return record, nil
}

func (s *Service) Delete(ctx context.Context, number string) (equipment.EquipmentRecord, error) {
	if err := ValidateNumber(number); err != nil {
		return equipment.EquipmentRecord{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.list.RemoveByNumber(number)
	if !ok {
		return equipment.EquipmentRecord{}, errors.New("record not found")
	}
	if err := s.persistLocked(ctx, audit.ActionDeleted, record.ID); err != nil {
		return equipment.EquipmentRecord{}, err
	}
	return record, nil
}

func (s *Service) UpdateStatus(ctx context.Context, number string, status string) error {
	if err := ValidateNumber(number); err != nil {
		return err
	}
	parsed, err := equipment.NormalizeStatus(status)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.list.UpdateStatus(number, parsed); err != nil {
		return err
	}
	node := s.list.FindByNumber(number)
	return s.persistLocked(ctx, audit.ActionStatusChanged, node.Value.ID)
}

func (s *Service) Query(ctx context.Context, borrower, status string) ([]equipment.EquipmentRecord, error) {
	if err := ValidateQuery(borrower, status); err != nil {
		return nil, err
	}
	s.mu.RLock()
	records := s.list.ToSlice()
	s.mu.RUnlock()
	if borrower != "" {
		records = equipment.FilterByBorrower(records, borrower)
	}
	if status != "" {
		parsed, _ := equipment.NormalizeStatus(status)
		records = equipment.FilterByStatus(records, parsed)
	}
	return records, nil
}

func (s *Service) Sort(ctx context.Context, mode string) ([]equipment.EquipmentRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []equipment.EquipmentRecord
	switch mode {
	case "number", "":
		result = equipment.SortByNumber(s.list.ToSlice())
	case "date":
		result = equipment.SortByBorrowDate(s.list.ToSlice())
	default:
		return nil, fmt.Errorf("unknown sort mode %q", mode)
	}
	s.list.Replace(result)
	if err := s.store.SaveRecords(ctx, result); err != nil {
		return nil, err
	}
	return cloneRecords(result), nil
}

func (s *Service) Snapshot() []equipment.EquipmentRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.list.ToSlice()
}

func (s *Service) persistLocked(ctx context.Context, action string, recordID int64) error {
	if err := s.store.SaveRecords(ctx, s.list.ToSlice()); err != nil {
		return err
	}
	entry := persistence.AuditEntry{Action: action, RecordID: recordID, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	entry.ID = time.Now().UnixNano()
	return s.store.SaveAudit(ctx, entry)
}

func nextID(records []equipment.EquipmentRecord) int64 {
	var largest int64
	for _, record := range records {
		if record.ID > largest {
			largest = record.ID
		}
	}
	return largest + 1
}

func cloneRecords(records []equipment.EquipmentRecord) []equipment.EquipmentRecord {
	clones := make([]equipment.EquipmentRecord, len(records))
	for index, record := range records {
		clones[index] = record.Clone()
	}
	return clones
}
