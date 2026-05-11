package aggregator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

// Template renders a new string field by evaluating a Go template against each
// JSON log entry. The result is stored under the configured output field.
type Template struct {
	field  string
	tmpl   *template.Template
	raw    string
	results []map[string]any
}

// NewTemplate creates a Template aggregator that evaluates tmplStr for each
// entry and stores the result in field. Returns an error if the template fails
// to parse.
func NewTemplate(field, tmplStr string) (*Template, error) {
	if field == "" {
		return nil, fmt.Errorf("template: field must not be empty")
	}
	if tmplStr == "" {
		return nil, fmt.Errorf("template: template string must not be empty")
	}
	t, err := template.New("logslice").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template: parse error: %w", err)
	}
	return &Template{field: field, tmpl: t, raw: tmplStr}, nil
}

// Record processes a single JSON log line, renders the template against the
// decoded object, and stores the rendered string under the configured field.
func (t *Template) Record(line []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, obj); err != nil {
		return fmt.Errorf("template: execute error: %w", err)
	}
	obj[t.field] = buf.String()
	t.results = append(t.results, obj)
	return nil
}

// Results returns all processed entries as JSON lines.
func (t *Template) Results() [][]byte {
	out := make([][]byte, 0, len(t.results))
	for _, obj := range t.results {
		b, err := json.Marshal(obj)
		if err != nil {
			continue
		}
		out = append(out, b)
	}
	return out
}

// Reset clears accumulated results.
func (t *Template) Reset() {
	t.results = t.results[:0]
}
