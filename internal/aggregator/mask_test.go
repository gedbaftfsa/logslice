package aggregator

import (
	"encoding/json"
	"testing"
)

func makeMaskEntry(t *testing.T, kv map[string]any) []byte {
	t.Helper()
	b, err := json.Marshal(kv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewMask_InvalidArgs(t *testing.T) {
	if _, err := NewMask("", `\d+`, ""); err == nil {
		t.Error("expected error for empty field")
	}
	if _, err := NewMask("email", "", ""); err == nil {
		t.Error("expected error for empty pattern")
	}
	if _, err := NewMask("email", `[invalid(", ""); err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestMask_InitialState(t *testing.T) {
	m, err := NewMask("token", `\w+`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap := m.Snapshot(); len(snap) != 0 {
		t.Errorf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestMask_Record_MasksWholeMatch(t *testing.T) {
	m, _ := NewMask("secret", `\d+`, "***")
	m.Record(makeMaskEntry(t, map[string]any{"secret": "abc123def456", "keep": "yes"}))
	snap := m.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	var row map[string]any
	if err := json.Unmarshal(snap[0], &row); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if row["secret"] != "abc***def***" {
		t.Errorf("expected abc***def***, got %v", row["secret"])
	}
	if row["keep"] != "yes" {
		t.Errorf("keep field should be unchanged")
	}
}

func TestMask_Record_MasksCapturingGroup(t *testing.T) {
	// pattern captures the local part of an email
	m, _ := NewMask("email", `(\w+)@`, "***")
	m.Record(makeMaskEntry(t, map[string]any{"email": "user@example.com"}))
	snap := m.Snapshot()
	var row map[string]any
	json.Unmarshal(snap[0], &row)
	if row["email"] != "***@example.com" {
		t.Errorf("expected ***@example.com, got %v", row["email"])
	}
}

func TestMask_Record_MissingFieldPassthrough(t *testing.T) {
	m, _ := NewMask("token", `\w+`, "***")
	m.Record(makeMaskEntry(t, map[string]any{"other": "value"}))
	snap := m.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry")
	}
	var row map[string]any
	json.Unmarshal(snap[0], &row)
	if row["other"] != "value" {
		t.Errorf("other field should be unchanged")
	}
}

func TestMask_Record_InvalidJSONSkipped(t *testing.T) {
	m, _ := NewMask("field", `\d+`, "***")
	m.Record([]byte(`not json`))
	if len(m.Snapshot()) != 0 {
		t.Error("invalid JSON should be skipped")
	}
}

func TestMask_Reset(t *testing.T) {
	m, _ := NewMask("x", `\d+`, "***")
	m.Record(makeMaskEntry(t, map[string]any{"x": "abc123"}))
	m.Reset()
	if len(m.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}
