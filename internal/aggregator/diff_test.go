package aggregator

import (
	"encoding/json"
	"fmt"
	"testing"
)

func makeDiffEntry(field string, value float64) []byte {
	b, _ := json.Marshal(map[string]interface{}{field: value})
	return b
}

func TestNewDiff_InitialState(t *testing.T) {
	d := NewDiff("latency")
	snap := d.Snapshot()
	if snap.Field != "latency" {
		t.Errorf("expected field latency, got %s", snap.Field)
	}
	if snap.Count != 0 {
		t.Errorf("expected 0 deltas, got %d", snap.Count)
	}
	if snap.Invalid != 0 {
		t.Errorf("expected 0 invalid, got %d", snap.Invalid)
	}
}

func TestDiff_Record_SingleEntry(t *testing.T) {
	d := NewDiff("bytes")
	d.Record(makeDiffEntry("bytes", 100))
	snap := d.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected 0 deltas after single entry, got %d", snap.Count)
	}
}

func TestDiff_Record_TwoEntries(t *testing.T) {
	d := NewDiff("bytes")
	d.Record(makeDiffEntry("bytes", 100))
	d.Record(makeDiffEntry("bytes", 150))
	snap := d.Snapshot()
	if snap.Count != 1 {
		t.Fatalf("expected 1 delta, got %d", snap.Count)
	}
	if snap.Deltas[0] != 50 {
		t.Errorf("expected delta 50, got %f", snap.Deltas[0])
	}
}

func TestDiff_Record_MultipleEntries(t *testing.T) {
	d := NewDiff("count")
	values := []float64{10, 20, 15, 30}
	expected := []float64{10, -5, 15}
	for _, v := range values {
		d.Record(makeDiffEntry("count", v))
	}
	snap := d.Snapshot()
	if snap.Count != len(expected) {
		t.Fatalf("expected %d deltas, got %d", len(expected), snap.Count)
	}
	for i, e := range expected {
		if snap.Deltas[i] != e {
			t.Errorf("delta[%d]: expected %f, got %f", i, e, snap.Deltas[i])
		}
	}
}

func TestDiff_Record_InvalidJSON(t *testing.T) {
	d := NewDiff("val")
	d.Record([]byte(`not json`))
	snap := d.Snapshot()
	if snap.Invalid != 1 {
		t.Errorf("expected 1 invalid, got %d", snap.Invalid)
	}
}

func TestDiff_Record_MissingField(t *testing.T) {
	d := NewDiff("missing")
	d.Record(makeDiffEntry("other", 42))
	d.Record(makeDiffEntry("other", 99))
	snap := d.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected 0 deltas for missing field, got %d", snap.Count)
	}
}

func TestDiff_Reset(t *testing.T) {
	d := NewDiff("x")
	d.Record(makeDiffEntry("x", 1))
	d.Record(makeDiffEntry("x", 5))
	d.Reset()
	snap := d.Snapshot()
	if snap.Count != 0 {
		t.Errorf("expected 0 after reset, got %d", snap.Count)
	}
}

func TestDiff_String(t *testing.T) {
	d := NewDiff("rps")
	d.Record(makeDiffEntry("rps", 10))
	d.Record(makeDiffEntry("rps", 20))
	s := d.String()
	expected := fmt.Sprintf("diff(rps): 1 deltas, 0 invalid")
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}
