# Histogram

The `Histogram` type tracks the statistical distribution of a numeric field
across structured JSON log entries.

## Usage

```go
h := aggregator.NewHistogram("duration_ms")

// Feed log lines
h.Record([]byte(`{"level":"info","duration_ms":42.5}`))
h.Record([]byte(`{"level":"info","duration_ms":120.0}`))

// Get a snapshot of current statistics
snap := h.Snapshot()
fmt.Printf("count=%d min=%.2f max=%.2f mean=%.2f p99=%.2f\n",
    snap.Count, snap.Min, snap.Max, snap.Mean, snap.P99)
```

## HistogramSnapshot fields

| Field   | Description                        |
|---------|------------------------------------|
| `field` | The tracked JSON field name        |
| `count` | Number of recorded values          |
| `min`   | Minimum observed value             |
| `max`   | Maximum observed value             |
| `mean`  | Arithmetic mean                    |
| `p50`   | 50th percentile (median)           |
| `p90`   | 90th percentile                    |
| `p99`   | 99th percentile                    |

## Notes

- Non-numeric values for the tracked field are silently skipped.
- Invalid JSON lines are ignored.
- `Reset()` clears all recorded values, useful for rolling windows.
- All methods are safe for concurrent use.
