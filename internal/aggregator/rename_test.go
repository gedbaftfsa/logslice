package aggregator

import (
	"encoding/json"
	"testing"
)

func makeRenameEntry(t *testing.T, fields map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewRename_InitialState(t *testing.T) {
	r := NewRename(map[string]string{"msg": "message"})
	if len(r.Results()) != 0 {
		t.Errorf("expected 0 results, got %d", len(r.Results()))
	}
}

func TestRename_Record_RenamesField(t *testing.T) {
	r := NewRename(map[string]string{"msg": "message"})
	entry := makeRenameEntry(t, map[string]interface{}{"msg": "hello", "level": "info"})
	if err := r.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := r.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	var out map[string]interface{}
	if err := json.Unmarshal(results[0], &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := out["msg"]; ok {
		t.Error("old key 'msg' should not exist")
	}
	if out["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", out["message"])
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestRename_Record_MissingField(t *testing.T) {
	r := NewRename(map[string]string{"missing": "renamed"})
	entry := makeRenameEntry(t, map[string]interface{}{"level": "warn"})
	if err := r.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := r.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	var out map[string]interface{}
	_ = json.Unmarshal(results[0], &out)
	if _, ok := out["renamed"]; ok {
		t.Error("renamed key should not appear when source key is absent")
	}
}

func TestRename_Record_InvalidJSON(t *testing.T) {
	r := NewRename(map[string]string{"a": "b"})
	if err := r.Record([]byte("not-json")); err == nil {
		t.Error("expected error for invalid JSON")
	}
	if len(r.Results()) != 0 {
		t.Error("expected 0 results after invalid JSON")
	}
}

func TestRename_Reset(t *testing.T) {
	r := NewRename(map[string]string{"a": "b"})
	entry := makeRenameEntry(t, map[string]interface{}{"a": "1"})
	_ = r.Record(entry)
	if len(r.Results()) != 1 {
		t.Fatalf("expected 1 result before reset")
	}
	r.Reset()
	if len(r.Results()) != 0 {
		t.Errorf("expected 0 results after reset, got %d", len(r.Results()))
	}
}

func TestRename_MultipleFields(t *testing.T) {
	r := NewRename(map[string]string{"ts": "timestamp", "msg": "message"})
	entry := makeRenameEntry(t, map[string]interface{}{"ts": "2024-01-01", "msg": "ok", "level": "info"})
	_ = r.Record(entry)
	results := r.Results()
	var out map[string]interface{}
	_ = json.Unmarshal(results[0], &out)
	if out["timestamp"] != "2024-01-01" {
		t.Errorf("expected timestamp=2024-01-01, got %v", out["timestamp"])
	}
	if out["message"] != "ok" {
		t.Errorf("expected message=ok, got %v", out["message"])
	}
	if _, ok := out["ts"]; ok {
		t.Error("old key 'ts' should not exist")
	}
}
