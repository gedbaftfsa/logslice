package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
)

func newWriter(buf *bytes.Buffer, format string) *output.Writer {
	w, _ := output.NewWriter(buf, format, nil)
	return w
}

func TestRun_NoFilter(t *testing.T) {
	input := `{"level":"info","msg":"hello"}
{"level":"error","msg":"oops"}
`
	var buf bytes.Buffer
	p := pipeline.New(strings.NewReader(input), nil, newWriter(&buf, "json"))
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 entries written, got %d", n)
	}
}

func TestRun_WithFilter(t *testing.T) {
	input := `{"level":"info","msg":"hello"}
{"level":"error","msg":"oops"}
{"level":"info","msg":"world"}
`
	var buf bytes.Buffer
	f, err := filter.Parse(`level="info"`)
	if err != nil {
		t.Fatalf("parse filter: %v", err)
	}
	p := pipeline.New(strings.NewReader(input), f, newWriter(&buf, "json"))
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 entries written, got %d", n)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	var buf bytes.Buffer
	p := pipeline.New(strings.NewReader(""), nil, newWriter(&buf, "json"))
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 entries, got %d", n)
	}
}

func TestRun_InvalidJSON(t *testing.T) {
	input := `not-json
`
	var buf bytes.Buffer
	p := pipeline.New(strings.NewReader(input), nil, newWriter(&buf, "json"))
	_, err := p.Run()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
