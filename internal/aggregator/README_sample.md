# aggregator/sample

Reservoir sampler that retains a statistically uniform random sample of up to **N** log entries from an unbounded stream.

## Algorithm

Uses [Algorithm R](https://en.wikipedia.org/wiki/Reservoir_sampling): the first N entries fill the reservoir; each subsequent entry replaces a random existing entry with probability N/count.

## Usage

```go
sampler := aggregator.NewSample(100)
for line := range src {
    sampler.Record(line)
}
for _, entry := range sampler.Snapshot() {
    fmt.Println(entry)
}
```

## CLI

```
cat app.log | logslice sample -n 50
```

Prints up to 50 randomly sampled log lines from stdin.

## Methods

| Method | Description |
|---|---|
| `NewSample(size int)` | Create sampler with given reservoir size |
| `Record(line []byte)` | Feed a raw JSON log line |
| `Snapshot() []map[string]any` | Return current sample entries |
| `Reset()` | Clear all entries |
| `MarshalJSON()` | Serialise state including count and entries |
