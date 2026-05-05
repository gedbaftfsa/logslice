package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format value.
// It returns an error if the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "json", "":
		return FormatJSON, nil
	case "pretty":
		return FormatPretty, nil
	case "text":
		return FormatText, nil
	default:
		return "", fmt.Errorf("output: unknown format %q (valid: json, pretty, text)", s)
	}
}

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// ValidFormats returns all supported format names.
func ValidFormats() []string {
	return []string{
		string(FormatJSON),
		string(FormatPretty),
		string(FormatText),
	}
}
