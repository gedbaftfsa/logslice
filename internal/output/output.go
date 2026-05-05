package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Format defines the output format for log entries.
type Format string

const (
	FormatJSON   Format = "json"
	FormatPretty Format = "pretty"
	FormatText   Format = "text"
)

// Writer writes log entries to an output destination.
type Writer struct {
	w      io.Writer
	format Format
	fields []string
}

// NewWriter creates a new Writer with the given format and optional field selection.
func NewWriter(w io.Writer, format Format, fields []string) *Writer {
	if w == nil {
		w = os.Stdout
	}
	return &Writer{w: w, format: format, fields: fields}
}

// Write outputs a single log entry (as a parsed map) to the writer.
func (wr *Writer) Write(entry map[string]any) error {
	if len(wr.fields) > 0 {
		entry = selectFields(entry, wr.fields)
	}
	switch wr.format {
	case FormatPretty:
		return wr.writePretty(entry)
	case FormatText:
		return wr.writeText(entry)
	default:
		return wr.writeJSON(entry)
	}
}

func (wr *Writer) writeJSON(entry map[string]any) error {
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(wr.w, string(b))
	return err
}

func (wr *Writer) writePretty(entry map[string]any) error {
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(wr.w, string(b))
	return err
}

func (wr *Writer) writeText(entry map[string]any) error {
	parts := make([]string, 0, len(entry))
	for k, v := range entry {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	_, err := fmt.Fprintln(wr.w, strings.Join(parts, " "))
	return err
}

func selectFields(entry map[string]any, fields []string) map[string]any {
	result := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := entry[f]; ok {
			result[f] = v
		}
	}
	return result
}
