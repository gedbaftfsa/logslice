package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/aggregator"
)

func TestRunPivot_Count(t *testing.T) {
	p, err := aggregator.NewPivot("level", "", "count")
	if err != nil {
		t.Fatalf("NewPivot: %v", err)
	}
	lines := []string{
		`{"level":"info","msg":"a"}`,
		`{"level":"info","msg":"b"}`,
		`{"level":"error","msg":"c"}`,
	}
	for _, l := range lines {
		p.Record([]byte(l))
	}
	snap := p.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("want 2 rows, got %d", len(snap))
	}
	if snap[1]["count"] != float64(2) {
		t.Errorf("info count: want 2, got %v", snap[1]["count"])
	}
}

func TestRunPivot_Sum(t *testing.T) {
	p, _ := aggregator.NewPivot("svc", "ms", "sum")
	entries := []map[string]interface{}{
		{"svc": "api", "ms": 10.0},
		{"svc": "api", "ms": 15.0},
		{"svc": "db", "ms": 3.0},
	}
	for _, e := range entries {
		b, _ := json.Marshal(e)
		p.Record(b)
	}
	snap := p.Snapshot()
	if snap[0]["sum"] != float64(25) {
		t.Errorf("api sum: want 25, got %v", snap[0]["sum"])
	}
}

func TestRunPivot_InvalidOp(t *testing.T) {
	_, err := aggregator.NewPivot("level", "dur", "variance")
	if err == nil {
		t.Error("expected error for unsupported op")
	}
}

func TestRunPivot_SnapshotJSON(t *testing.T) {
	p, _ := aggregator.NewPivot("env", "req", "max")
	p.Record([]byte(`{"env":"prod","req":99.0}`))
	p.Record([]byte(`{"env":"prod","req":50.0}`))

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, row := range p.Snapshot() {
		if err := enc.Encode(row); err != nil {
			t.Fatalf("encode: %v", err)
		}
	}
	out := buf.String()
	if !strings.Contains(out, "99") {
		t.Errorf("expected max=99 in output, got: %s", out)
	}
}
