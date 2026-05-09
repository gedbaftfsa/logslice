package aggregator

import (
	"encoding/json"
	"testing"
)

func makeTruncateEntry(t *testing.T, m map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewTruncate_InvalidArgs(t *testing.T) {
	_, err := NewTruncate("", 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
	_, err = NewTruncate("msg", 0)
	if err == nil {
		t.Fatal("expected error for zero maxLen")
	}
	_, err = NewTruncate("msg", -5)
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestTruncate_ShortStringUnchanged(t *testing.T) {
	tr, _ := NewTruncate("msg", 20)
	entry := makeTruncateEntry(t, map[string]interface{}{"msg": "hello", "level": "info"})
	if err := tr.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}
	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 result, got %d", len(snap))
	}
	var out map[string]interface{}
	json.Unmarshal(snap[0], &out)
	if out["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", out["msg"])
	}
}

func TestTruncate_LongStringTrimmed(t *testing.T) {
	tr, _ := NewTruncate("msg", 5)
	entry := makeTruncateEntry(t, map[string]interface{}{"msg": "hello world"})
	tr.Record(entry)
	snap := tr.Snapshot()
	var out map[string]interface{}
	json.Unmarshal(snap[0], &out)
	if out["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", out["msg"])
	}
}

func TestTruncate_MissingFieldPassthrough(t *testing.T) {
	tr, _ := NewTruncate("msg", 5)
	entry := makeTruncateEntry(t, map[string]interface{}{"level": "warn"})
	tr.Record(entry)
	snap := tr.Snapshot()
	var out map[string]interface{}
	json.Unmarshal(snap[0], &out)
	if _, ok := out["msg"]; ok {
		t.Error("expected msg field to be absent")
	}
	if out["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", out["level"])
	}
}

func TestTruncate_NonStringFieldUnchanged(t *testing.T) {
	tr, _ := NewTruncate("count", 3)
	entry := makeTruncateEntry(t, map[string]interface{}{"count": 42})
	tr.Record(entry)
	snap := tr.Snapshot()
	var out map[string]interface{}
	json.Unmarshal(snap[0], &out)
	if out["count"] != float64(42) {
		t.Errorf("expected count=42, got %v", out["count"])
	}
}

func TestTruncate_InvalidJSON(t *testing.T) {
	tr, _ := NewTruncate("msg", 5)
	if err := tr.Record([]byte("not json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestTruncate_Reset(t *testing.T) {
	tr, _ := NewTruncate("msg", 5)
	tr.Record(makeTruncateEntry(t, map[string]interface{}{"msg": "hello"}))
	tr.Reset()
	if len(tr.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}
