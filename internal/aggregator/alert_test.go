package aggregator

import (
	"encoding/json"
	"testing"
)

func makeAlertEntry(t *testing.T, fields map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewAlert_InitialState(t *testing.T) {
	a := NewAlert(AlertCondition{Field: "latency", Operator: "gt", Threshold: 100})
	snap := a.Snapshot()
	var m map[string]interface{}
	if err := json.Unmarshal(snap, &m); err != nil {
		t.Fatalf("unmarshal snapshot: %v", err)
	}
	if m["triggered"].(float64) != 0 {
		t.Errorf("expected 0 triggered, got %v", m["triggered"])
	}
	if m["firing"].(bool) {
		t.Error("expected firing=false initially")
	}
}

func TestAlert_Record_TriggersGT(t *testing.T) {
	a := NewAlert(AlertCondition{Field: "latency", Operator: "gt", Threshold: 100})
	a.Record(makeAlertEntry(t, map[string]interface{}{"latency": 200.0}))
	a.Record(makeAlertEntry(t, map[string]interface{}{"latency": 50.0}))
	var m map[string]interface{}
	json.Unmarshal(a.Snapshot(), &m)
	if m["triggered"].(float64) != 1 {
		t.Errorf("expected 1 triggered, got %v", m["triggered"])
	}
	if !m["firing"].(bool) {
		t.Error("expected firing=true")
	}
}

func TestAlert_Record_TriggersLT(t *testing.T) {
	a := NewAlert(AlertCondition{Field: "score", Operator: "lt", Threshold: 0.5})
	a.Record(makeAlertEntry(t, map[string]interface{}{"score": 0.1}))
	a.Record(makeAlertEntry(t, map[string]interface{}{"score": 0.9}))
	var m map[string]interface{}
	json.Unmarshal(a.Snapshot(), &m)
	if m["triggered"].(float64) != 1 {
		t.Errorf("expected 1 triggered, got %v", m["triggered"])
	}
}

func TestAlert_Record_InvalidJSON(t *testing.T) {
	a := NewAlert(AlertCondition{Field: "x", Operator: "gt", Threshold: 1})
	a.Record([]byte("not-json"))
	var m map[string]interface{}
	json.Unmarshal(a.Snapshot(), &m)
	if m["total"].(float64) != 0 {
		t.Errorf("expected total=0 for invalid JSON")
	}
}

func TestAlert_Record_MissingField(t *testing.T) {
	a := NewAlert(AlertCondition{Field: "latency", Operator: "gt", Threshold: 10})
	a.Record(makeAlertEntry(t, map[string]interface{}{"other": 999.0}))
	var m map[string]interface{}
	json.Unmarshal(a.Snapshot(), &m)
	if m["triggered"].(float64) != 0 {
		t.Errorf("expected 0 triggered for missing field")
	}
}

func TestAlert_Reset(t *testing.T) {
	a := NewAlert(AlertCondition{Field: "latency", Operator: "gt", Threshold: 10})
	a.Record(makeAlertEntry(t, map[string]interface{}{"latency": 999.0}))
	a.Reset()
	var m map[string]interface{}
	json.Unmarshal(a.Snapshot(), &m)
	if m["triggered"].(float64) != 0 {
		t.Errorf("expected 0 after reset")
	}
}
