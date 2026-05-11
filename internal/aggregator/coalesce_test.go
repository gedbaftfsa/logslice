package aggregator

import (
	"encoding/json"
	"testing"
)

func makeCoalesceEntry(t *testing.T, kv map[string]any) []byte {
	t.Helper()
	b, err := json.Marshal(kv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewCoalesce_InvalidArgs(t *testing.T) {
	_, err := NewCoalesce([]string{"only"}, "dest")
	if err == nil {
		t.Fatal("expected error for single field")
	}
	_, err = NewCoalesce([]string{"a", "b"}, "")
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestCoalesce_PicksFirstNonEmpty(t *testing.T) {
	c, _ := NewCoalesce([]string{"msg", "message", "text"}, "log")

	entry := makeCoalesceEntry(t, map[string]any{"message": "hello", "level": "info"})
	if err := c.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	snap := c.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 result, got %d", len(snap))
	}

	var out map[string]any
	if err := json.Unmarshal(snap[0], &out); err != nil {
		t.Fatal(err)
	}
	if out["log"] != "hello" {
		t.Errorf("expected log=hello, got %v", out["log"])
	}
}

func TestCoalesce_SkipsEmptyString(t *testing.T) {
	c, _ := NewCoalesce([]string{"msg", "message"}, "log")

	entry := makeCoalesceEntry(t, map[string]any{"msg": "", "message": "world"})
	_ = c.Record(entry)

	var out map[string]any
	_ = json.Unmarshal(c.Snapshot()[0], &out)
	if out["log"] != "world" {
		t.Errorf("expected log=world, got %v", out["log"])
	}
}

func TestCoalesce_NoMatchPassthrough(t *testing.T) {
	c, _ := NewCoalesce([]string{"msg", "message"}, "log")

	entry := makeCoalesceEntry(t, map[string]any{"level": "warn"})
	_ = c.Record(entry)

	var out map[string]any
	_ = json.Unmarshal(c.Snapshot()[0], &out)
	if _, ok := out["log"]; ok {
		t.Error("expected no log field when no source matched")
	}
	if out["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", out["level"])
	}
}

func TestCoalesce_InvalidJSON(t *testing.T) {
	c, _ := NewCoalesce([]string{"a", "b"}, "out")
	if err := c.Record([]byte("not-json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCoalesce_Reset(t *testing.T) {
	c, _ := NewCoalesce([]string{"a", "b"}, "out")
	_ = c.Record(makeCoalesceEntry(t, map[string]any{"a": "v"}))
	c.Reset()
	if len(c.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}
