package aggregator

import (
	"encoding/json"
	"strings"
	"testing"
)

func makeExtractEntry(kv map[string]any) string {
	b, _ := json.Marshal(kv)
	return string(b)
}

func TestNewExtract_InvalidArgs(t *testing.T) {
	_, err := NewExtract("", "dst", `(\d+)`)
	if err == nil {
		t.Fatal("expected error for empty srcField")
	}
	_, err = NewExtract("src", "", `(\d+)`)
	if err == nil {
		t.Fatal("expected error for empty dstField")
	}
	_, err = NewExtract("src", "dst", `[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
	_, err = NewExtract("src", "dst", `nodgroup`)
	if err == nil {
		t.Fatal("expected error for pattern without capture group")
	}
}

func TestExtract_Record_ExtractsMatch(t *testing.T) {
	ex, err := NewExtract("msg", "code", `ERROR-(\d+)`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := makeExtractEntry(map[string]any{"msg": "prefix ERROR-42 suffix", "level": "error"})
	_ = ex.Record(line)

	snap := ex.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(snap[0]), &m)
	if m["code"] != "42" {
		t.Errorf("expected code=42, got %v", m["code"])
	}
}

func TestExtract_Record_MissingField(t *testing.T) {
	ex, _ := NewExtract("msg", "code", `(\d+)`)
	line := makeExtractEntry(map[string]any{"level": "info"})
	_ = ex.Record(line)

	snap := ex.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	var m map[string]any
	_ = json.Unmarshal([]byte(snap[0]), &m)
	if _, ok := m["code"]; ok {
		t.Error("expected code field to be absent")
	}
}

func TestExtract_Record_InvalidJSON(t *testing.T) {
	ex, _ := NewExtract("msg", "code", `(\d+)`)
	_ = ex.Record("not json")
	if len(ex.Snapshot()) != 0 {
		t.Error("expected no entries for invalid JSON")
	}
}

func TestExtract_Record_NoMatch(t *testing.T) {
	ex, _ := NewExtract("msg", "code", `ERROR-(\d+)`)
	line := makeExtractEntry(map[string]any{"msg": "all good"})
	_ = ex.Record(line)

	snap := ex.Snapshot()
	var m map[string]any
	_ = json.Unmarshal([]byte(snap[0]), &m)
	if _, ok := m["code"]; ok {
		t.Error("expected no code field when pattern does not match")
	}
}

func TestExtract_Reset(t *testing.T) {
	ex, _ := NewExtract("msg", "code", `(\d+)`)
	_ = ex.Record(makeExtractEntry(map[string]any{"msg": "val 7"}))
	ex.Reset()
	if len(ex.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}

func TestExtract_MultipleEntries(t *testing.T) {
	ex, _ := NewExtract("path", "filename", `/([^/]+)$`)
	lines := []string{
		makeExtractEntry(map[string]any{"path": "/var/log/app.log"}),
		makeExtractEntry(map[string]any{"path": "/etc/config.yaml"}),
	}
	for _, l := range lines {
		_ = ex.Record(l)
	}
	snap := ex.Snapshot()
	expected := []string{"app.log", "config.yaml"}
	for i, s := range snap {
		var m map[string]any
		_ = json.Unmarshal([]byte(s), &m)
		if !strings.EqualFold(fmt.Sprintf("%v", m["filename"]), expected[i]) {
			t.Errorf("entry %d: expected filename=%s, got %v", i, expected[i], m["filename"])
		}
	}
}
