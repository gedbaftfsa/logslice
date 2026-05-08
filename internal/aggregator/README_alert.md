# Alert

The `Alert` aggregator monitors a numeric field in JSON log entries and counts how many entries exceed (or fall below, or equal) a configured threshold.

## Usage

```go
cond := aggregator.AlertCondition{
    Field:     "latency_ms",
    Operator:  "gt",   // "gt", "lt", or "eq"
    Threshold: 500,
}
a := aggregator.NewAlert(cond)

for _, line := range logLines {
    a.Record(line)
}

fmt.Println(a.String())
// alert[latency_ms gt 500.00]: 3/10 triggered
```

## Snapshot Output

```json
{
  "field": "latency_ms",
  "operator": "gt",
  "threshold": 500,
  "triggered": 3,
  "total": 10,
  "firing": true
}
```

## Operators

| Operator | Meaning          |
|----------|------------------|
| `gt`     | greater than     |
| `lt`     | less than        |
| `eq`     | equal to         |

## Notes

- Non-numeric field values are silently skipped.
- Invalid JSON entries do not increment any counters.
- `Reset()` clears all counters without changing the condition.
