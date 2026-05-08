package aggregator

import (
	"encoding/json"
	"fmt"
	"testing"
)

func makeSampleEntry(val string) []byte {
	return []byte(fmt.Sprintf(`{"msg":%q}`, val))
}

func TestNewSample_InitialState(t *testing.T) {
	s := NewSample(5)
	if got := s.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(got))
	}
}

func TestSample_Record_ValidJSON(t *testing.T) {
	s := NewSample(3)
	for i := 0; i < 3; i++ {
		s.Record(makeSampleEntry(fmt.Sprintf("msg%d", i)))
	}
	if got := len(s.Snapshot()); got != 3 {
		t.Fatalf("expected 3 entries, got %d", got)
	}
}

func TestSample_Record_InvalidJSON(t *testing.T) {
	s := NewSample(5)
	s.Record([]byte(`not-json`))
	if got := len(s.Snapshot()); got != 0 {
		t.Fatalf("expected 0 entries after invalid JSON, got %d", got)
	}
}

func TestSample_ReservoirCap(t *testing.T) {
	s := NewSample(5)
	for i := 0; i < 100; i++ {
		s.Record(makeSampleEntry(fmt.Sprintf("entry%d", i)))
	}
	if got := len(s.Snapshot()); got != 5 {
		t.Fatalf("expected reservoir capped at 5, got %d", got)
	}
}

func TestSample_Reset(t *testing.T) {
	s := NewSample(5)
	for i := 0; i < 5; i++ {
		s.Record(makeSampleEntry(fmt.Sprintf("e%d", i)))
	}
	s.Reset()
	if got := len(s.Snapshot()); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestSample_MarshalJSON(t *testing.T) {
	s := NewSample(3)
	s.Record(makeSampleEntry("hello"))
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if _, ok := out["entries"]; !ok {
		t.Fatal("expected 'entries' key in JSON output")
	}
	if _, ok := out["sample_size"]; !ok {
		t.Fatal("expected 'sample_size' key in JSON output")
	}
}
