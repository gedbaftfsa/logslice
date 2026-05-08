package aggregator

import (
	"fmt"
	"testing"
	"time"
)

func makeDedupeEntry(field, value string) []byte {
	return []byte(fmt.Sprintf(`{"%s":"%s"}`, field, value))
}

func TestNewDedupe_InitialState(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	if d.Dropped != 0 {
		t.Fatalf("expected 0 dropped, got %d", d.Dropped)
	}
}

func TestDedupe_FirstOccurrenceNotDuplicate(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	entry := makeDedupeEntry("msg", "hello")
	dup, err := d.IsDuplicate(entry)
	if err != nil {
		t.Fatal(err)
	}
	if dup {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestDedupe_SecondOccurrenceIsDuplicate(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	entry := makeDedupeEntry("msg", "hello")
	_, _ = d.IsDuplicate(entry)
	dup, err := d.IsDuplicate(entry)
	if err != nil {
		t.Fatal(err)
	}
	if !dup {
		t.Fatal("second occurrence within TTL should be a duplicate")
	}
	if d.Dropped != 1 {
		t.Fatalf("expected 1 dropped, got %d", d.Dropped)
	}
}

func TestDedupe_DifferentValuesNotDuplicate(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	_, _ = d.IsDuplicate(makeDedupeEntry("msg", "hello"))
	dup, err := d.IsDuplicate(makeDedupeEntry("msg", "world"))
	if err != nil {
		t.Fatal(err)
	}
	if dup {
		t.Fatal("different values should not be duplicates")
	}
}

func TestDedupe_MissingFieldNotDuplicate(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	dup, err := d.IsDuplicate([]byte(`{"other":"value"}`))
	if err != nil {
		t.Fatal(err)
	}
	if dup {
		t.Fatal("missing field should not be treated as duplicate")
	}
}

func TestDedupe_InvalidJSON(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	_, err := d.IsDuplicate([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDedupe_Reset(t *testing.T) {
	d := NewDedupe("msg", time.Second)
	entry := makeDedupeEntry("msg", "hello")
	_, _ = d.IsDuplicate(entry)
	_, _ = d.IsDuplicate(entry)
	d.Reset()
	if d.Dropped != 0 {
		t.Fatalf("expected 0 dropped after reset, got %d", d.Dropped)
	}
	dup, _ := d.IsDuplicate(entry)
	if dup {
		t.Fatal("after reset, entry should not be duplicate")
	}
}

func TestDedupe_ExpiredEntryNotDuplicate(t *testing.T) {
	d := NewDedupe("msg", 10*time.Millisecond)
	entry := makeDedupeEntry("msg", "hello")
	_, _ = d.IsDuplicate(entry)
	time.Sleep(20 * time.Millisecond)
	dup, err := d.IsDuplicate(entry)
	if err != nil {
		t.Fatal(err)
	}
	if dup {
		t.Fatal("expired entry should not be a duplicate")
	}
}
