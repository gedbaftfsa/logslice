# aggregator/counter

The `Counter` type tallies how often each distinct value appears for a chosen
field across a stream of JSON log entries.

## Usage

```go
c := aggregator.NewCounter("level")

for entry := range logStream {
    c.Record(entry)
}

counts := c.Snapshot()
fmt.Println(counts) // map[info:42 error:3 warn:7]
```

## API

| Function / Method | Description |
|---|---|
| `NewCounter(field string) *Counter` | Create a counter for the named field |
| `Record(entry []byte)` | Parse entry and increment the matching value's count |
| `Snapshot() map[string]int64` | Return a copy of current counts (thread-safe) |
| `Reset()` | Clear all counts |
| `MarshalJSON()` | Serialise as `{"field":"...","counts":{...}}` |

## Notes

- Only string-valued fields are counted; numeric and boolean values are ignored.
- Entries with missing fields or invalid JSON are silently skipped.
- All methods are safe for concurrent use.
