# Limit Aggregator

The `Limit` aggregator passes through at most **N** log entries and then signals
that no more input is needed. It is useful for previewing a stream or capping
output in scripts.

## Usage

```go
lim, err := aggregator.NewLimit(100)
if err != nil {
    log.Fatal(err)
}

for _, line := range lines {
    more := lim.Record(line)
    if !more {
        break
    }
}

for _, entry := range lim.Entries() {
    // process entry
}
```

## CLI flag

```
logslice --limit 50 app.log
```

## Behaviour

| Scenario | Result |
|---|---|
| `n <= 0` | `NewLimit` returns an error |
| invalid JSON line | skipped, does not consume a slot |
| count < n | `Record` returns `true` (keep feeding) |
| count == n | `Record` returns `false`; `Done()` is `true` |

## Snapshot

`Snapshot()` returns a JSON object:

```json
{"max": 50, "count": 50, "entries": [...]}
```
