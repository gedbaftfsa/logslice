package aggregator

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func makeEntry(ts time.Time, msg string) json.RawMessage {
	s := fmt.Sprintf(`{"time":%q,"msg":%q}`, ts.Format(time.RFC3339), msg)
	return json.RawMessage(s)
}

func TestWindow_EmptyInitially(t *testing.T) {
	w := NewWindow(time.Minute)
	if w.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", w.Len())
	}
}

func TestWindow_AddsEntry(t *testing.T) {
	w := NewWindow(time.Minute)
	now := time.Now()
	w.Add(makeEntry(now, "hello"), "time")
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", w.Len())
	}
}

func TestWindow_EvictsOldEntries(t *testing.T) {
	w := NewWindow(30 * time.Second)
	now := time.Now()
	old := now.Add(-60 * time.Second)
	w.Add(makeEntry(old, "old"), "time")
	w.Add(makeEntry(now, "new"), "time")
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", w.Len())
	}
	entries := w.Entries()
	var m map[string]string
	if err := json.Unmarshal(entries[0], &m); err != nil {
		t.Fatal(err)
	}
	if m["msg"] != "new" {
		t.Fatalf("expected 'new', got %q", m["msg"])
	}
}

func TestWindow_KeepsAllWithinDuration(t *testing.T) {
	w := NewWindow(time.Minute)
	now := time.Now()
	for i := 0; i < 5; i++ {
		w.Add(makeEntry(now.Add(time.Duration(i)*time.Second), fmt.Sprintf("msg%d", i)), "time")
	}
	if w.Len() != 5 {
		t.Fatalf("expected 5 entries, got %d", w.Len())
	}
}

func TestWindow_FallbackOnMissingField(t *testing.T) {
	w := NewWindow(time.Hour)
	w.Add(json.RawMessage(`{"msg":"no-ts"}`), "time")
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", w.Len())
	}
}

func TestWindow_FallbackOnEmptyField(t *testing.T) {
	w := NewWindow(time.Hour)
	w.Add(json.RawMessage(`{"msg":"no-ts"}`), "")
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", w.Len())
	}
}
