package audit

import (
	"fmt"
	"sort"
	"strings"

	"equipmentlending/internal/persistence"
)

const (
	ActionCreated       = "created"
	ActionDeleted       = "deleted"
	ActionStatusChanged = "status_changed"
	ActionExported      = "exported"
)

func Classify(action string) string {
	switch strings.ToLower(strings.TrimSpace(action)) {
	case ActionCreated:
		return "record_lifecycle"
	case ActionDeleted:
		return "record_lifecycle"
	case ActionStatusChanged:
		return "record_state"
	case ActionExported:
		return "data_delivery"
	default:
		return "other"
	}
}

func Summary(entries []persistence.AuditEntry) string {
	if len(entries) == 0 {
		return "no audit events"
	}
	counts := make(map[string]int)
	for _, entry := range entries {
		counts[Classify(entry.Action)]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, counts[key]))
	}
	return strings.Join(parts, ", ")
}

func Timeline(entries []persistence.AuditEntry) []string {
	ordered := append([]persistence.AuditEntry(nil), entries...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].CreatedAt < ordered[j].CreatedAt })
	result := make([]string, 0, len(ordered))
	for _, entry := range ordered {
		result = append(result, fmt.Sprintf("%s: record=%d action=%s", entry.CreatedAt, entry.RecordID, entry.Action))
	}
	return result
}

func IsMutation(action string) bool {
	return action == ActionCreated || action == ActionDeleted || action == ActionStatusChanged
}
