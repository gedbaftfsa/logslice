package aggregator_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/aggregator"
)

// chanSource wraps a slice of lines as an aggregator.Source.
type chanSource struct{ lines []string }

func (c *chanSource) Lines() <-chan string {
	ch := make(chan string, len(c.lines))
	for _, l := range c.lines {
		ch <- l
	}
	close(ch)
	return ch
}

func collect(ch <-chan aggregator.Entry) []aggregator.Entry {
	var out []aggregator.Entry
	for e := range ch {
		out = append(out, e)
	}
	return out
}

func TestMerge_SingleSource(t *testing.T) {
	src := &chanSource{lines: []string{`{"level":"info","msg":"hello"}`}}
	ch := aggregator.Merge(map[string]aggregator.Source{"app": src})
	entries := collect(ch)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Source != "app" {
		t.Errorf("expected source 'app', got %q", entries[0].Source)
	}
	if entries[0].Fields["msg"] != "hello" {
		t.Errorf("unexpected fields: %v", entries[0].Fields)
	}
}

func TestMerge_MultipleSources(t *testing.T) {
	sources := map[string]aggregator.Source{
		"svc-a": &chanSource{lines: []string{`{"msg":"a1"}`, `{"msg":"a2"}`}},
		"svc-b": &chanSource{lines: []string{`{"msg":"b1"}`}},
	}
	entries := collect(aggregator.Merge(sources))
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	labels := make([]string, len(entries))
	for i, e := range entries {
		labels[i] = e.Source
	}
	sort.Strings(labels)
	if strings.Join(labels, ",") != "svc-a,svc-a,svc-b" {
		t.Errorf("unexpected labels: %v", labels)
	}
}

func TestMerge_SkipsInvalidJSON(t *testing.T) {
	src := &chanSource{lines: []string{"not json", `{"ok":true}`}}
	entries := collect(aggregator.Merge(map[string]aggregator.Source{"x": src}))
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestMerge_Empty(t *testing.T) {
	entries := collect(aggregator.Merge(map[string]aggregator.Source{}))
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}
