# Coalesce

The `Coalesce` processor picks the **first non-empty string value** from an
ordered list of source fields and writes it to a destination field.

This is useful when different log producers use different field names for the
same concept (e.g. `msg`, `message`, `text`) and you want a single canonical
field downstream.

## Usage

```go
c, err := aggregator.NewCoalesce(
    []string{"msg", "message", "text"}, // tried in order
    "log",                               // destination field
)
```

## Behaviour

| Situation | Result |
|-----------|--------|
| First field present and non-empty | Written to `dest`; entry passed through |
| First field missing or empty | Next field tried |
| No field matches | Entry passed through unchanged |
| Invalid JSON | Error returned; entry skipped |

## Methods

- `Record(line []byte) error` — process one JSON log line.
- `Snapshot() [][]byte` — return all processed entries as JSON lines.
- `Reset()` — clear accumulated results.

## Example

Input:
```json
{"message": "hello", "level": "info"}
```

Output with `fields=["msg","message"]`, `dest="log"`:
```json
{"message": "hello", "level": "info", "log": "hello"}
```
