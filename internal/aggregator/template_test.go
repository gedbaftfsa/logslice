package aggregator

import (
	"encoding/json"
	"testing"
)

func makeTemplateEntry(t *testing.T, fields map[string]any) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewTemplate_InvalidArgs(t *testing.T) {
	if _, err := NewTemplate("", "{{.level}}"); err == nil {
		t.Fatal("expected error for empty field")
	}
	if _, err := NewTemplate("msg", ""); err == nil {
		t.Fatal("expected error for empty template")
	}
	if _, err := NewTemplate("msg", "{{.unclosed"); err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestTemplate_Record_RendersField(t *testing.T) {
	tmpl, err := NewTemplate("summary", "[{{.level}}] {{.msg}}")
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}
	entry := makeTemplateEntry(t, map[string]any{"level": "error", "msg": "disk full"})
	if err := tmpl.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}
	results := tmpl.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	var obj map[string]any
	if err := json.Unmarshal(results[0], &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := obj["summary"]; got != "[error] disk full" {
		t.Errorf("summary = %q, want %q", got, "[error] disk full")
	}
}

func TestTemplate_Record_MissingFieldZero(t *testing.T) {
	tmpl, err := NewTemplate("summary", "{{.level}}-{{.missing}}")
	if err != nil {
		t.Fatalf("NewTemplate: %v", err)
	}
	entry := makeTemplateEntry(t, map[string]any{"level": "info"})
	if err := tmpl.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}
	results := tmpl.Results()
	var obj map[string]any
	json.Unmarshal(results[0], &obj)
	if got := obj["summary"]; got != "info-<no value>" {
		t.Errorf("summary = %q", got)
	}
}

func TestTemplate_Record_InvalidJSON(t *testing.T) {
	tmpl, _ := NewTemplate("summary", "{{.level}}")
	if err := tmpl.Record([]byte("not-json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestTemplate_Reset(t *testing.T) {
	tmpl, _ := NewTemplate("out", "{{.x}}")
	tmpl.Record(makeTemplateEntry(t, map[string]any{"x": "1"}))
	tmpl.Reset()
	if len(tmpl.Results()) != 0 {
		t.Fatal("expected empty results after Reset")
	}
}
