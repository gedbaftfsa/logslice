# TimeBucket

The `TimeBucket` aggregator groups JSON log entries into fixed-duration time
buckets and counts how many entries fall into each bucket.

## Usage

```go
tb := aggregator.NewTimeBucket("ts", time.Minute)

for _, line := range lines {
    tb.Record(line)
}

snap := tb.Snapshot() // map[time.Time]int64
```

## Fields

| Parameter  | Description                                          |
|------------|------------------------------------------------------|
| `field`    | JSON key that holds the timestamp value              |
| `interval` | Bucket width (e.g. `time.Minute`, `5*time.Minute`)   |

## Supported timestamp formats

- **RFC 3339** string — `"2024-01-01T12:00:00Z"`
- **Unix epoch** number — `1704110400`

Entries whose timestamp field is missing or unparseable are silently skipped.

## CLI

```
cat app.log | logslice timebucket --field ts --interval 5m
```

Outputs a JSON object mapping each bucket's start time (RFC 3339) to its count:

```json
{
  "2024-01-01T12:00:00Z": 42,
  "2024-01-01T12:05:00Z": 17
}
```

## Methods

- `Record(line []byte)` — parse and bucket one log line
- `Snapshot() map[time.Time]int64` — return a copy of current counts
- `Reset()` — clear all buckets
- `MarshalJSON() ([]byte, error)` — serialise as a JSON object
