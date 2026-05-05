package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func entry() map[string]any {
	return map[string]any{
		"level":   "info",
		"message": "hello world",
		"ts":      1234567890,
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON, nil)
	if err := w.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := strings.TrimSpace(buf.String())
	var out map[string]any
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestWrite_Pretty(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatPretty, nil)
	if err := w.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\n") {
		t.Error("pretty output should contain newlines")
	}
}

func TestWrite_Text(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatText, nil)
	if err := w.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "level=info") {
		t.Errorf("text output missing level=info, got: %s", output)
	}
}

func TestWrite_FieldSelection(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON, []string{"level"})
	if err := w.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := strings.TrimSpace(buf.String())
	var out map[string]any
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := out["message"]; ok {
		t.Error("field 'message' should have been excluded")
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestWrite_FieldSelection_Missing(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON, []string{"nonexistent"})
	if err := w.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := strings.TrimSpace(buf.String())
	var out map[string]any
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty object, got %v", out)
	}
}
