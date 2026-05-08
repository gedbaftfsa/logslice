# Dedupe

The `Dedupe` aggregator suppresses duplicate log entries based on a specified field value within a configurable time-to-live (TTL) window.

## Usage

```go
d := aggregator.NewDedupe("request_id", 5*time.Second)

for entry := range source {
    dup, err := d.IsDuplicate(entry)
    if err != nil || dup {
        continue // skip duplicates and invalid JSON
    }
    // process unique entry
}
```

## Behaviour

- The first occurrence of a field value is always passed through.
- Subsequent occurrences of the same value within the TTL are suppressed and counted in `Dropped`.
- Once the TTL expires the value is evicted and the next occurrence is treated as new.
- Entries missing the target field are never treated as duplicates.
- Invalid JSON returns an error and is not recorded.

## Methods

| Method | Description |
|---|---|
| `NewDedupe(field, ttl)` | Create a new Dedupe for the given field and TTL |
| `IsDuplicate(entry)` | Returns `(true, nil)` if the entry is a duplicate within TTL |
| `Reset()` | Clear all seen entries and reset the dropped counter |

## Fields

- `Dropped int` — total number of entries suppressed since creation or last reset.
