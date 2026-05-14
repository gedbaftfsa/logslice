# Sort

The `Sort` aggregator buffers JSON log entries and emits them in sorted order by a specified field.

## Usage

```go
s, err := aggregator.NewSort("status_code", "asc")
if err != nil {
    log.Fatal(err)
}

for _, line := range lines {
    s.Record(line)
}

sorted, err := s.Flush()
if err != nil {
    log.Fatal(err)
}
for _, b := range sorted {
    fmt.Println(string(b))
}
```

## Parameters

| Parameter | Description |
|-----------|-------------|
| `field`   | JSON field name to sort by |
| `order`   | `"asc"` for ascending, `"desc"` for descending |

## Supported Field Types

- **Numeric** (`float64` after JSON decode)
- **String**

Entries whose field is missing or of an unsupported type are treated as equal and preserve their relative order (stable sort).

## Notes

- All entries are buffered in memory until `Flush()` is called.
- Call `Reset()` to clear the buffer between runs.
- Invalid JSON lines return an error from `Record()` and are not buffered.
