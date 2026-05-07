package aggregator

import (
	"encoding/json"
	"testing"
)

func makeGroupEntry(t *testing.T, fields map[string]string) []byte {
	t.Helper()
	m := make(map[string]interface{})
	for k, v := range fields {
		m[k] = v
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewGroupBy_InitialState(t *testing.T) {
	g := NewGroupBy("level")
	snap := g.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %v", snap)
	}
}

func TestGroupBy_Record_StringField(t *testing.T) {
	g := NewGroupBy("level")
	g.Record(makeGroupEntry(t, map[string]string{"level": "info"}))
	g.Record(makeGroupEntry(t, map[string]string{"level": "info"}))
	g.Record(makeGroupEntry(t, map[string]string{"level": "error"}))
	snap := g.Snapshot()
	if snap["info"] != 2 {
		t.Errorf("expected info=2, got %d", snap["info"])
	}
	if snap["error"] != 1 {
		t.Errorf("expected error=1, got %d", snap["error"])
	}
}

func TestGroupBy_Record_MissingField(t *testing.T) {
	g := NewGroupBy("level")
	g.Record(makeGroupEntry(t, map[string]string{"msg": "hello"}))
	if len(g.Snapshot()) != 0 {
		t.Error("expected no counts for missing field")
	}
}

func TestGroupBy_Record_InvalidJSON(t *testing.T) {
	g := NewGroupBy("level")
	g.Record([]byte("not-json"))
	if len(g.Snapshot()) != 0 {
		t.Error("expected no counts for invalid JSON")
	}
}

func TestGroupBy_Reset(t *testing.T) {
	g := NewGroupBy("level")
	g.Record(makeGroupEntry(t, map[string]string{"level": "warn"}))
	g.Reset()
	if len(g.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}

func TestGroupBy_MarshalJSON(t *testing.T) {
	g := NewGroupBy("level")
	g.Record(makeGroupEntry(t, map[string]string{"level": "debug"}))
	b, err := json.Marshal(g)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var out map[string]int
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if out["debug"] != 1 {
		t.Errorf("expected debug=1, got %d", out["debug"])
	}
}
