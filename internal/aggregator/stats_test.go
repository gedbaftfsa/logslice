package aggregator

import (
	"encoding/json"
	"sync"
	"testing"
)

func TestNewStats_InitialValues(t *testing.T) {
	s := NewStats()
	if s.Total != 0 {
		t.Errorf("expected Total 0, got %d", s.Total)
	}
	if s.Invalid != 0 {
		t.Errorf("expected Invalid 0, got %d", s.Invalid)
	}
	if len(s.BySource) != 0 {
		t.Errorf("expected empty BySource, got %v", s.BySource)
	}
}

func TestStats_Record(t *testing.T) {
	s := NewStats()
	s.Record("app.log")
	s.Record("app.log")
	s.Record("svc.log")

	if s.Total != 3 {
		t.Errorf("expected Total 3, got %d", s.Total)
	}
	if s.BySource["app.log"] != 2 {
		t.Errorf("expected app.log count 2, got %d", s.BySource["app.log"])
	}
	if s.BySource["svc.log"] != 1 {
		t.Errorf("expected svc.log count 1, got %d", s.BySource["svc.log"])
	}
}

func TestStats_RecordInvalid(t *testing.T) {
	s := NewStats()
	s.RecordInvalid()
	s.RecordInvalid()

	if s.Invalid != 2 {
		t.Errorf("expected Invalid 2, got %d", s.Invalid)
	}
	if s.Total != 0 {
		t.Errorf("expected Total 0, got %d", s.Total)
	}
}

func TestStats_Snapshot(t *testing.T) {
	s := NewStats()
	s.Record("a")
	s.RecordInvalid()

	snap := s.Snapshot()
	if snap["total"].(int64) != 1 {
		t.Errorf("snapshot total mismatch")
	}
	if snap["invalid"].(int64) != 1 {
		t.Errorf("snapshot invalid mismatch")
	}
}

func TestStats_MarshalJSON(t *testing.T) {
	s := NewStats()
	s.Record("x")

	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if out["total"] == nil {
		t.Error("expected 'total' key in JSON output")
	}
}

func TestStats_ConcurrentRecord(t *testing.T) {
	s := NewStats()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Record("concurrent")
		}()
	}
	wg.Wait()
	if s.Total != 100 {
		t.Errorf("expected Total 100 after concurrent writes, got %d", s.Total)
	}
}
