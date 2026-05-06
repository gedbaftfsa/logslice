package aggregator

import (
	"fmt"
	"testing"
	"time"
)

func TestNewRate_InitialState(t *testing.T) {
	r := NewRate(5 * time.Second)
	if r.Count() != 0 {
		t.Errorf("expected 0 initial count, got %d", r.Count())
	}
	if r.PerSecond() != 0 {
		t.Errorf("expected 0 initial rate, got %f", r.PerSecond())
	}
}

func TestRate_Record_ValidJSON(t *testing.T) {
	r := NewRate(5 * time.Second)
	for i := 0; i < 5; i++ {
		r.Record([]byte(fmt.Sprintf(`{"msg":"event %d"}`, i)))
	}
	if r.Count() != 5 {
		t.Errorf("expected 5 events, got %d", r.Count())
	}
}

func TestRate_Record_InvalidJSON(t *testing.T) {
	r := NewRate(5 * time.Second)
	r.Record([]byte(`not-json`))
	if r.Count() != 0 {
		t.Errorf("expected 0 events for invalid JSON, got %d", r.Count())
	}
}

func TestRate_PerSecond(t *testing.T) {
	r := NewRate(2 * time.Second)
	for i := 0; i < 4; i++ {
		r.Record([]byte(`{"msg":"x"}`))
	}
	rate := r.PerSecond()
	// 4 events over 2s window = 2.0/s
	if rate < 1.5 || rate > 2.5 {
		t.Errorf("expected ~2.0 events/sec, got %f", rate)
	}
}

func TestRate_Reset(t *testing.T) {
	r := NewRate(5 * time.Second)
	r.Record([]byte(`{"msg":"x"}`))
	r.Record([]byte(`{"msg":"y"}`))
	r.Reset()
	if r.Count() != 0 {
		t.Errorf("expected 0 after reset, got %d", r.Count())
	}
}

func TestRate_Eviction(t *testing.T) {
	r := NewRate(100 * time.Millisecond)
	r.Record([]byte(`{"msg":"old"}`))
	time.Sleep(150 * time.Millisecond)
	// After sleeping, the old entry should be evicted
	if r.Count() != 0 {
		t.Errorf("expected 0 after eviction, got %d", r.Count())
	}
}
