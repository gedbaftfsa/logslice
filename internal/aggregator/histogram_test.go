package aggregator

import (
	"encoding/json"
	"testing"
)

func makeEntry(t *testing.T, kv map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(kv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewHistogram_InitialState(t *testing.T) {
	h := NewHistogram("duration_ms")
	snap := h.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected count 0, got %d", snap.Count)
	}
	if snap.Field != "duration_ms" {
		t.Errorf("expected field duration_ms, got %s", snap.Field)
	}
}

func TestHistogram_Record_ValidField(t *testing.T) {
	h := NewHistogram("latency")
	for _, v := range []float64{10, 20, 30, 40, 50} {
		h.Record(makeEntry(t, map[string]interface{}{"latency": v}))
	}
	snap := h.Snapshot()
	if snap.Count != 5 {
		t.Errorf("expected count 5, got %d", snap.Count)
	}
	if snap.Min != 10 {
		t.Errorf("expected min 10, got %f", snap.Min)
	}
	if snap.Max != 50 {
		t.Errorf("expected max 50, got %f", snap.Max)
	}
	if snap.Mean != 30 {
		t.Errorf("expected mean 30, got %f", snap.Mean)
	}
}

func TestHistogram_Record_MissingField(t *testing.T) {
	h := NewHistogram("latency")
	h.Record(makeEntry(t, map[string]interface{}{"other": 42.0}))
	snap := h.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected count 0, got %d", snap.Count)
	}
}

func TestHistogram_Record_InvalidJSON(t *testing.T) {
	h := NewHistogram("latency")
	h.Record([]byte("not-json"))
	snap := h.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected count 0 for invalid JSON, got %d", snap.Count)
	}
}

func TestHistogram_Percentiles(t *testing.T) {
	h := NewHistogram("ms")
	for i := 1; i <= 100; i++ {
		h.Record(makeEntry(t, map[string]interface{}{"ms": float64(i)}))
	}
	snap := h.Snapshot()
	if snap.P50 != 50 {
		t.Errorf("expected p50=50, got %f", snap.P50)
	}
	if snap.P90 != 90 {
		t.Errorf("expected p90=90, got %f", snap.P90)
	}
	if snap.P99 != 99 {
		t.Errorf("expected p99=99, got %f", snap.P99)
	}
}

func TestHistogram_Reset(t *testing.T) {
	h := NewHistogram("ms")
	h.Record(makeEntry(t, map[string]interface{}{"ms": 100.0}))
	h.Reset()
	snap := h.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected count 0 after reset, got %d", snap.Count)
	}
}
