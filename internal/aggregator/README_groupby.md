# GroupBy Aggregator

The `GroupBy` type counts log entries grouped by the distinct values of a
named JSON field.

## Usage

```go
g := aggregator.NewGroupBy("level")

// Feed raw JSON lines
g.Record([]byte(`{"level":"info","msg":"started"}`))
g.Record([]byte(`{"level":"error","msg":"failed"}`))
g.Record([]byte(`{"level":"info","msg":"done"}`))

// Read current counts
snap := g.Snapshot()
// snap == map[string]int{"info": 2, "error": 1}

// Serialise to JSON
b, _ := json.Marshal(g)
// b == {"error":1,"info":2}

// Reset for a new window
g.Reset()
```

## Behaviour

- Lines that are not valid JSON are silently skipped.
- Lines where the target field is absent are silently skipped.
- Only string field values are supported; non-string values are skipped.
- `Snapshot` and `Record` are safe for concurrent use.

## CLI

The `groupby` sub-command in `internal/cli/groupby_flag.go` wraps this
aggregator and prints the result as indented JSON to stdout:

```
cat app.log | logslice groupby -field level
```
