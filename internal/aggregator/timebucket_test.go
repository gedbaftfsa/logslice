package aggregator

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func makeTBEntry(ts string) []byte {
	return []byte(fmt.Sprintf(`{"ts":%q,"msg":"hello"}`, ts))
}

func TestNewTimeBucket_InitialState(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	snap := tb.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestTimeBucket_Record_ValidRFC3339(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	tb.Record(makeTBEntry("2024-01-01T12:00:30Z"))
	tb.Record(makeTBEntry("2024-01-01T12:00:45Z"))
	snap := tb.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	for _, count := range snap {
		if count != 2 {
			t.Errorf("expected count 2, got %d", count)
		}
	}
}

func TestTimeBucket_Record_MultipleBuckets(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	tb.Record(makeTBEntry("2024-01-01T12:00:00Z"))
	tb.Record(makeTBEntry("2024-01-01T12:01:00Z"))
	snap := tb.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(snap))
	}
}

func TestTimeBucket_Record_UnixTimestamp(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	now := time.Now().Truncate(time.Minute).Unix()
	entry := []byte(fmt.Sprintf(`{"ts":%d}`, now))
	tb.Record(entry)
	snap := tb.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
}

func TestTimeBucket_Record_MissingField(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	tb.Record([]byte(`{"msg":"no timestamp"}`))
	if len(tb.Snapshot()) != 0 {
		t.Error("expected no buckets for missing field")
	}
}

func TestTimeBucket_Record_InvalidJSON(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	tb.Record([]byte(`not json`))
	if len(tb.Snapshot()) != 0 {
		t.Error("expected no buckets for invalid JSON")
	}
}

func TestTimeBucket_Reset(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	tb.Record(makeTBEntry("2024-01-01T12:00:00Z"))
	tb.Reset()
	if len(tb.Snapshot()) != 0 {
		t.Error("expected empty snapshot after reset")
	}
}

func TestTimeBucket_MarshalJSON(t *testing.T) {
	tb := NewTimeBucket("ts", time.Minute)
	tb.Record(makeTBEntry("2024-01-01T12:00:00Z"))
	b, err := json.Marshal(tb)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var out map[string]int64
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}
