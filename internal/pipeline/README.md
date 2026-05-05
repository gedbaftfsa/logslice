# pipeline

The `pipeline` package wires together a log source, an optional filter, and an output writer into a single processing unit.

## Usage

```go
import (
    "os"
    "github.com/yourorg/logslice/internal/filter"
    "github.com/yourorg/logslice/internal/output"
    "github.com/yourorg/logslice/internal/pipeline"
)

f, _ := filter.Parse(`level="error"`)
w, _ := output.NewWriter(os.Stdout, "pretty", nil)
p := pipeline.New(os.Stdin, f, w)
n, err := p.Run()
```

## Behaviour

- Reads newline-delimited JSON objects from any `io.Reader`.
- Skips entries that do not match the provided filter (pass `nil` to disable filtering).
- Writes matching entries through the configured `output.Writer`.
- Returns the count of written entries and the first non-EOF error encountered.
