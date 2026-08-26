package audit

import (
	"strings"
	"testing"

	"equipmentlending/internal/persistence"
)

func TestAuditSummary(t *testing.T) {
	entries := []persistence.AuditEntry{{Action: ActionCreated}, {Action: ActionStatusChanged}, {Action: ActionCreated}}
	text := Summary(entries)
	if !strings.Contains(text, "record_lifecycle=2") || !strings.Contains(text, "record_state=1") {
		t.Fatalf("unexpected summary %s", text)
	}
	if !IsMutation(ActionDeleted) || Classify("unknown") != "other" {
		t.Fatal("audit classification failed")
	}
}
