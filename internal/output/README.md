# output

The `output` package handles formatting and writing of filtered log entries to a destination.

## Overview

After log entries are read from a source and matched against a filter, the `output.Writer` serializes them in the desired format.

## Formats

| Format   | Description                              |
|----------|------------------------------------------|
| `json`   | Compact single-line JSON (default)       |
| `pretty` | Indented multi-line JSON                 |
| `text`   | Space-separated `key=value` pairs        |

## Field Selection

You can limit the output to specific fields by passing a non-empty `fields` slice to `NewWriter`. Only those keys will appear in the output.

## Usage

```go
import "github.com/yourorg/logslice/internal/output"

w := output.NewWriter(os.Stdout, output.FormatJSON, []string{"level", "message"})

entry := map[string]any{
    "level":   "error",
    "message": "disk full",
    "host":    "web-01",
}

if err := w.Write(entry); err != nil {
    log.Fatal(err)
}
// Output: {"level":"error","message":"disk full"}
```
