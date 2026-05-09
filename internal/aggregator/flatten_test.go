package aggregator

import (
	"encoding/json"
	"testing"
)

func makeFlattenEntry(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewFlatten_InitialState(t *testing.T) {
	f := NewFlatten("")
	if got := f.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(got))
	}
}

func TestFlatten_FlatObject(t *testing.T) {
	f := NewFlatten("")
	entry := makeFlattenEntry(t, map[string]any{"level": "info", "msg": "hello"})
	if err := f.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := f.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0]["level"] != "info" || snap[0]["msg"] != "hello" {
		t.Errorf("unexpected entry: %v", snap[0])
	}
}

func TestFlatten_NestedObject(t *testing.T) {
	f := NewFlatten("")
	entry := makeFlattenEntry(t, map[string]any{
		"http": map[string]any{"method": "GET", "status": float64(200)},
	})
	if err := f.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := f.Snapshot()
	if snap[0]["http.method"] != "GET" {
		t.Errorf("expected http.method=GET, got %v", snap[0]["http.method"])
	}
	if snap[0]["http.status"] != float64(200) {
		t.Errorf("expected http.status=200, got %v", snap[0]["http.status"])
	}
}

func TestFlatten_DeeplyNested(t *testing.T) {
	f := NewFlatten("")
	entry := makeFlattenEntry(t, map[string]any{
		"a": map[string]any{"b": map[string]any{"c": "deep"}},
	})
	if err := f.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := f.Snapshot()
	if snap[0]["a.b.c"] != "deep" {
		t.Errorf("expected a.b.c=deep, got %v", snap[0]["a.b.c"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	f := NewFlatten("log")
	entry := makeFlattenEntry(t, map[string]any{"level": "warn"})
	if err := f.Record(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := f.Snapshot()
	if snap[0]["log.level"] != "warn" {
		t.Errorf("expected log.level=warn, got %v", snap[0]["log.level"])
	}
}

func TestFlatten_InvalidJSON(t *testing.T) {
	f := NewFlatten("")
	if err := f.Record([]byte("not-json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if len(f.Snapshot()) != 0 {
		t.Fatal("expected no entries after invalid JSON")
	}
}

func TestFlatten_Reset(t *testing.T) {
	f := NewFlatten("")
	_ = f.Record(makeFlattenEntry(t, map[string]any{"x": 1}))
	f.Reset()
	if len(f.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot after Reset")
	}
}
