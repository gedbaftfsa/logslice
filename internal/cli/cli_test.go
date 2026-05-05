package cli

import (
	"testing"
)

func TestParseArgs_Defaults(t *testing.T) {
	cfg, err := parseArgs([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Filter != "" {
		t.Errorf("expected empty filter, got %q", cfg.Filter)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Format)
	}
	if len(cfg.Fields) != 0 {
		t.Errorf("expected no fields, got %v", cfg.Fields)
	}
	if len(cfg.Files) != 0 {
		t.Errorf("expected no files, got %v", cfg.Files)
	}
}

func TestParseArgs_WithFlags(t *testing.T) {
	cfg, err := parseArgs([]string{
		"-filter", "level=error",
		"-format", "pretty",
		"-fields", "level, msg, ts",
		"file1.log", "file2.log",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Filter != "level=error" {
		t.Errorf("expected filter 'level=error', got %q", cfg.Filter)
	}
	if cfg.Format != "pretty" {
		t.Errorf("expected format 'pretty', got %q", cfg.Format)
	}
	expectedFields := []string{"level", "msg", "ts"}
	if len(cfg.Fields) != len(expectedFields) {
		t.Fatalf("expected %d fields, got %d", len(expectedFields), len(cfg.Fields))
	}
	for i, f := range expectedFields {
		if cfg.Fields[i] != f {
			t.Errorf("field[%d]: expected %q, got %q", i, f, cfg.Fields[i])
		}
	}
	if len(cfg.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(cfg.Files))
	}
}

func TestParseArgs_InvalidFlag(t *testing.T) {
	_, err := parseArgs([]string{"-unknown"})
	if err == nil {
		t.Error("expected error for unknown flag, got nil")
	}
}

func TestParseArgs_EmptyFields(t *testing.T) {
	cfg, err := parseArgs([]string{"-fields", ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 0 {
		t.Errorf("expected no fields for empty string, got %v", cfg.Fields)
	}
}
