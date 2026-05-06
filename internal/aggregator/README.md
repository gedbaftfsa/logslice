# aggregator

The `aggregator` package fans-in log entries from multiple named sources into a
single channel, making it easy for the pipeline to process logs from several
files or streams simultaneously.

## Key types

| Symbol | Description |
|--------|-------------|
| `Source` | Interface satisfied by any type that exposes a `Lines() <-chan string` method. |
| `Entry` | A parsed log record carrying its origin label, decoded fields, and the original raw line. |

## Functions

### `Merge(sources map[string]Source) <-chan Entry`

Starts one goroutine per source, reads all lines concurrently, parses each line
as JSON, and emits `Entry` values on the returned channel. Non-JSON lines are
silently dropped. The channel is closed when every source is exhausted.

### `InjectSourceField(fields map[string]interface{}, label string) map[string]interface{}`

Returns a shallow copy of `fields` with the `_source` key set to `label`.
Useful when the output formatter needs to include provenance in each record.

### `FormatLabel(label string) string`

Returns `[label]` — a short prefix used by the text output formatter.

## Usage

```go
entries := aggregator.Merge(map[string]aggregator.Source{
    "api":    apiFileSource,
    "worker": workerTailSource,
})
for e := range entries {
    fmt.Println(e.Source, e.Fields)
}
```
