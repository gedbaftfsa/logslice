package source

import (
	"os"
	"strings"
	"testing"
)

func collectLines(s Source) []string {
	var result []string
	for line := range s.Lines() {
		result = append(result, line)
	}
	return result
}

func TestNewReaderSource_Empty(t *testing.T) {
	src := NewReaderSource(strings.NewReader(""))
	lines := collectLines(src)
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestNewReaderSource_SingleLine(t *testing.T) {
	input := `{"level":"info","msg":"hello"}`
	src := NewReaderSource(strings.NewReader(input))
	lines := collectLines(src)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != input {
		t.Errorf("expected %q, got %q", input, lines[0])
	}
}

func TestNewReaderSource_MultipleLines(t *testing.T) {
	input := "{\"a\":1}\n{\"b\":2}\n{\"c\":3}"
	src := NewReaderSource(strings.NewReader(input))
	lines := collectLines(src)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestNewReaderSource_SkipsBlankLines(t *testing.T) {
	input := "{\"a\":1}\n\n{\"b\":2}\n"
	src := NewReaderSource(strings.NewReader(input))
	lines := collectLines(src)
	if len(lines) != 2 {
		t.Fatalf("expected 2 non-blank lines, got %d", len(lines))
	}
}

func TestNewFileSource_NotFound(t *testing.T) {
	_, err := NewFileSource("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestNewFileSource_Valid(t *testing.T) {
	f, err := os.CreateTemp("", "logslice-test-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString("{\"level\":\"debug\"}\n{\"level\":\"error\"}\n")
	f.Close()

	src, err := NewFileSource(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := collectLines(src)
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
