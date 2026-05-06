package aggregator

import (
	"encoding/json"
	"testing"
)

func TestNewCounter_InitialState(t *testing.T) {
	c := NewCounter("level")
	if c.field != "level" {
		t.Fatalf("expected field 'level', got %q", c.field)
	}
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected empty counts initially")
	}
}

func TestCounter_Record_StringField(t *testing.T) {
	c := NewCounter("level")
	entries := []string{
		`{"level":"info","msg":"a"}`,
		`{"level":"info","msg":"b"}`,
		`{"level":"error","msg":"c"}`,
	}
	for _, e := range entries {
		c.Record([]byte(e))
	}
	snap := c.Snapshot()
	if snap["info"] != 2 {
		t.Errorf("expected info=2, got %d", snap["info"])
	}
	if snap["error"] != 1 {
		t.Errorf("expected error=1, got %d", snap["error"])
	}
}

func TestCounter_Record_MissingField(t *testing.T) {
	c := NewCounter("level")
	c.Record([]byte(`{"msg":"no level here"}`))
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected no counts for missing field")
	}
}

func TestCounter_Record_InvalidJSON(t *testing.T) {
	c := NewCounter("level")
	c.Record([]byte(`not json`))
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected no counts for invalid JSON")
	}
}

func TestCounter_Reset(t *testing.T) {
	c := NewCounter("level")
	c.Record([]byte(`{"level":"info"}`))
	c.Reset()
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected empty counts after reset")
	}
}

func TestCounter_MarshalJSON(t *testing.T) {
	c := NewCounter("level")
	c.Record([]byte(`{"level":"warn"}`))
	c.Record([]byte(`{"level":"warn"}`))

	data, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	var out struct {
		Field  string           `json:"field"`
		Counts map[string]int64 `json:"counts"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if out.Field != "level" {
		t.Errorf("expected field 'level', got %q", out.Field)
	}
	if out.Counts["warn"] != 2 {
		t.Errorf("expected warn=2, got %d", out.Counts["warn"])
	}
}
