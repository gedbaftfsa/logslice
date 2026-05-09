package aggregator

import (
	"encoding/json"
	"testing"
)

func makeRedactEntry(t *testing.T, fields map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewRedact_NoFields(t *testing.T) {
	_, err := NewRedact(nil, "")
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestRedact_DefaultPlaceholder(t *testing.T) {
	r, err := NewRedact([]string{"password"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.placeholder != "***" {
		t.Errorf("expected default placeholder ***, got %s", r.placeholder)
	}
}

func TestRedact_Record_RedactsField(t *testing.T) {
	r, _ := NewRedact([]string{"token"}, "[REDACTED]")
	entry := makeRedactEntry(t, map[string]interface{}{"user": "alice", "token": "secret123"})
	out, err := r.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(out, &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj["token"] != "[REDACTED]" {
		t.Errorf("expected token to be redacted, got %v", obj["token"])
	}
	if obj["user"] != "alice" {
		t.Errorf("user field should be unchanged, got %v", obj["user"])
	}
}

func TestRedact_Record_MissingField(t *testing.T) {
	r, _ := NewRedact([]string{"password"}, "")
	entry := makeRedactEntry(t, map[string]interface{}{"msg": "hello"})
	out, err := r.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(entry) {
		t.Errorf("expected unchanged output, got %s", out)
	}
	if r.count != 0 {
		t.Errorf("expected count 0, got %d", r.count)
	}
}

func TestRedact_Record_InvalidJSON(t *testing.T) {
	r, _ := NewRedact([]string{"x"}, "")
	out, err := r.Record([]byte("not-json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "not-json" {
		t.Errorf("expected passthrough for invalid JSON")
	}
}

func TestRedact_Snapshot(t *testing.T) {
	r, _ := NewRedact([]string{"ssn", "dob"}, "MASKED")
	entry := makeRedactEntry(t, map[string]interface{}{"ssn": "123-45-6789", "dob": "1990-01-01"})
	r.Record(entry)

	snap := r.Snapshot()
	var obj map[string]interface{}
	if err := json.Unmarshal(snap, &obj); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	if obj["lines_redacted"].(float64) != 1 {
		t.Errorf("expected lines_redacted=1, got %v", obj["lines_redacted"])
	}
}

func TestRedact_Reset(t *testing.T) {
	r, _ := NewRedact([]string{"key"}, "")
	entry := makeRedactEntry(t, map[string]interface{}{"key": "val"})
	r.Record(entry)
	r.Reset()
	if r.count != 0 {
		t.Errorf("expected count 0 after reset, got %d", r.count)
	}
}
