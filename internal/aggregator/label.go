package aggregator

import "fmt"

// InjectSourceField returns a copy of fields with the "_source" key set to
// label. Existing fields are not mutated.
func InjectSourceField(fields map[string]interface{}, label string) map[string]interface{} {
	out := make(map[string]interface{}, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	out["_source"] = label
	return out
}

// FormatLabel returns a display-friendly prefix for a source label used when
// rendering entries in text mode.
func FormatLabel(label string) string {
	return fmt.Sprintf("[%s]", label)
}
