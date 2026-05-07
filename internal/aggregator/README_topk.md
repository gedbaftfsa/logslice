# TopK

The `TopK` type tracks the most frequent string values for a specified JSON field across a stream of log lines.

## Usage

```go
tk := aggregator.NewTopK("level", 5)

for _, line := range lines {
    tk.Record(line)
}

for _, entry := range tk.Top() {
    fmt.Printf("%s: %d\n", entry.Value, entry.Count)
}
```

## API

### `NewTopK(field string, k int) *TopK`

Creates a new TopK tracker. If `k <= 0`, it defaults to `10`.

### `Record(line []byte) error`

Parses a JSON log line and increments the count for the value of the tracked field. Non-string values and missing fields are silently skipped.

### `Top() []TopKEntry`

Returns up to `k` entries sorted by count descending. Ties are broken alphabetically.

### `Reset()`

Clears all counts.

## Thread Safety

All methods are safe for concurrent use.
