# Rename

The `Rename` transformer renames fields in structured JSON log entries.

## Usage

```go
rename := aggregator.NewRename(map[string]string{
    "msg":   "message",
    "ts":    "timestamp",
    "level": "severity",
})

for _, line := range lines {
    _ = rename.Record(line)
}

for _, result := range rename.Results() {
    fmt.Println(string(result))
}
```

## Behaviour

- Each call to `Record` parses the JSON line and applies all configured renames.
- If a source key does not exist in the entry, it is silently skipped.
- The original key is removed; the value is stored under the new key.
- `Results()` returns all transformed entries as marshalled JSON lines.
- `Reset()` clears all stored results without modifying the mapping.
- Invalid JSON lines return an error and are not stored.

## Thread Safety

`Record`, `Results`, and `Reset` are protected by a mutex and safe for concurrent use.
