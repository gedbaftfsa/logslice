package aggregator

import (
	"encoding/json"
	"testing"
)

func makeLimitEntry(t *testing.T, kv map[string]any) []byte {
	t.Helper()
	b, err := json.Marshal(kv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewLimit_InvalidN(t *testing.T) {
	_, err := NewLimit(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
	_, err = NewLimit(-5)
	if err == nil {
		t.Fatal("expected error for n=-5")
	}
}

func TestNewLimit_InitialState(t *testing.T) {
	l, err := NewLimit(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Done() {
		t.Error("should not be done initially")
	}
	if len(l.Entries()) != 0 {
		t.Error("entries should be empty initially")
	}
}

func TestLimit_Record_UnderLimit(t *testing.T) {
	l, _ := NewLimit(3)
	entry := makeLimitEntry(t, map[string]any{"msg": "hello"})
	more := l.Record(entry)
	if !more {
		t.Error("expected more=true when under limit")
	}
	if len(l.Entries()) != 1 {
		t.Errorf("expected 1 entry, got %d", len(l.Entries()))
	}
}

func TestLimit_Record_AtLimit(t *testing.T) {
	l, _ := NewLimit(2)
	entry := makeLimitEntry(t, map[string]any{"msg": "a"})
	l.Record(entry)
	more := l.Record(entry)
	if more {
		t.Error("expected more=false when limit reached")
	}
	if !l.Done() {
		t.Error("expected Done()=true")
	}
}

func TestLimit_Record_InvalidJSON(t *testing.T) {
	l, _ := NewLimit(3)
	more := l.Record([]byte("not-json"))
	if !more {
		t.Error("invalid JSON should not consume a slot")
	}
	if len(l.Entries()) != 0 {
		t.Error("invalid JSON should not be stored")
	}
}

func TestLimit_Reset(t *testing.T) {
	l, _ := NewLimit(2)
	entry := makeLimitEntry(t, map[string]any{"x": 1})
	l.Record(entry)
	l.Record(entry)
	l.Reset()
	if l.Done() {
		t.Error("should not be done after reset")
	}
	if len(l.Entries()) != 0 {
		t.Error("entries should be empty after reset")
	}
}

func TestLimit_Snapshot(t *testing.T) {
	l, _ := NewLimit(5)
	entry := makeLimitEntry(t, map[string]any{"level": "info"})
	l.Record(entry)
	b, err := l.Snapshot()
	if err != nil {
		t.Fatalf("snapshot error: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	if out["max"].(float64) != 5 {
		t.Errorf("expected max=5, got %v", out["max"])
	}
	if out["count"].(float64) != 1 {
		t.Errorf("expected count=1, got %v", out["count"])
	}
}
