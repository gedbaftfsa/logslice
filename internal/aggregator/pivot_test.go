package aggregator

import (
	"encoding/json"
	"testing"
)

func makePivotEntry(t *testing.T, fields map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewPivot_InvalidOp(t *testing.T) {
	_, err := NewPivot("level", "duration", "median")
	if err == nil {
		t.Fatal("expected error for invalid op")
	}
}

func TestPivot_Count(t *testing.T) {
	p, _ := NewPivot("level", "", "count")
	p.Record(makePivotEntry(t, map[string]interface{}{"level": "info"}))
	p.Record(makePivotEntry(t, map[string]interface{}{"level": "info"}))
	p.Record(makePivotEntry(t, map[string]interface{}{"level": "error"}))

	snap := p.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(snap))
	}
	if snap[0]["level"] != "error" || snap[0]["count"] != float64(1) {
		t.Errorf("unexpected error bucket: %v", snap[0])
	}
	if snap[1]["level"] != "info" || snap[1]["count"] != float64(2) {
		t.Errorf("unexpected info bucket: %v", snap[1])
	}
}

func TestPivot_Sum(t *testing.T) {
	p, _ := NewPivot("service", "latency", "sum")
	p.Record(makePivotEntry(t, map[string]interface{}{"service": "api", "latency": 10.0}))
	p.Record(makePivotEntry(t, map[string]interface{}{"service": "api", "latency": 20.0}))
	p.Record(makePivotEntry(t, map[string]interface{}{"service": "db", "latency": 5.0}))

	snap := p.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(snap))
	}
	if snap[0]["sum"] != float64(30) {
		t.Errorf("api sum: want 30, got %v", snap[0]["sum"])
	}
	if snap[1]["sum"] != float64(5) {
		t.Errorf("db sum: want 5, got %v", snap[1]["sum"])
	}
}

func TestPivot_Avg(t *testing.T) {
	p, _ := NewPivot("region", "score", "avg")
	p.Record(makePivotEntry(t, map[string]interface{}{"region": "us", "score": 80.0}))
	p.Record(makePivotEntry(t, map[string]interface{}{"region": "us", "score": 100.0}))

	snap := p.Snapshot()
	if snap[0]["avg"] != float64(90) {
		t.Errorf("avg: want 90, got %v", snap[0]["avg"])
	}
}

func TestPivot_MissingKeyField(t *testing.T) {
	p, _ := NewPivot("level", "dur", "sum")
	p.Record(makePivotEntry(t, map[string]interface{}{"msg": "hello", "dur": 5.0}))
	if len(p.Snapshot()) != 0 {
		t.Error("expected empty snapshot when key field missing")
	}
}

func TestPivot_InvalidJSON(t *testing.T) {
	p, _ := NewPivot("level", "dur", "count")
	p.Record([]byte(`not-json`))
	if len(p.Snapshot()) != 0 {
		t.Error("expected empty snapshot on invalid JSON")
	}
}

func TestPivot_Reset(t *testing.T) {
	p, _ := NewPivot("level", "", "count")
	p.Record(makePivotEntry(t, map[string]interface{}{"level": "info"}))
	p.Reset()
	if len(p.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}
