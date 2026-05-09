# Flatten

The `Flatten` transformer expands nested JSON objects into a flat map using dot-notation keys.

## Use case

Structured logs often contain deeply nested fields (e.g. `http.request.method`). Flatten makes every field addressable at the top level, simplifying downstream filtering and field selection.

## Constructor

```go
f := aggregator.NewFlatten(prefix string) *Flatten
```

- `prefix` — optional string prepended to every key (`""` for none).

## Methods

| Method | Description |
|---|---|
| `Record(line []byte) error` | Ingest a raw JSON line and store the flattened result. |
| `Snapshot() []map[string]any` | Return all flattened entries recorded so far. |
| `Reset()` | Clear all stored entries. |

## Example

Input:
```json
{"http":{"method":"GET","status":200},"level":"info"}
```

Output after `Flatten` with prefix `""`:
```json
{"http.method":"GET","http.status":200,"level":"info"}
```

With prefix `"log"`:
```json
{"log.http.method":"GET","log.http.status":200,"log.level":"info"}
```

## Notes

- Only `map[string]any` values are recursed into; arrays and scalar values are kept as-is.
- The original entry is never mutated.
- Safe for concurrent use.
