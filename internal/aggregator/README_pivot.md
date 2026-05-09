# Pivot Aggregator

The `Pivot` aggregator groups JSON log entries by a **key field** and computes a
summary statistic over a **value field** for each group.

## Supported Operations

| Op      | Description                        |
|---------|------------------------------------|
| `count` | Number of entries per key          |
| `sum`   | Sum of the value field per key     |
| `avg`   | Average of the value field per key |
| `min`   | Minimum value per key              |
| `max`   | Maximum value per key              |

## Usage

```go
p, err := aggregator.NewPivot("service", "latency_ms", "avg")
if err != nil {
    log.Fatal(err)
}
for _, line := range lines {
    p.Record(line)
}
for _, row := range p.Snapshot() {
    fmt.Println(row)
}
```

## Output

Each row in the snapshot is a `map[string]interface{}` with the key field and
the aggregated value under the operation name:

```json
{"service": "api", "avg": 42.5}
{"service": "db",  "avg": 8.1}
```

Rows are sorted alphabetically by key.

## Notes

- Entries missing the key field or (for non-count ops) the value field are skipped.
- Non-numeric value fields are skipped.
- Call `Reset()` to clear accumulated state between intervals.
