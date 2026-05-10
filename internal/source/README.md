# source

The `source` package provides readers for structured JSON log lines from multiple input origins.

## Types

### `Source`

An interface representing any log input stream. Implementations must provide a `Lines()` method returning a channel of raw JSON strings.

## Constructors

### `NewReaderSource(r io.Reader) Source`

Creates a source from any `io.Reader`. Reads line by line.

### `NewFileSource(path string) (Source, error)`

Opens a file at the given path and returns a source for it. Returns an error if the file cannot be opened.

### `NewStdinSource() Source`

Returns a source that reads from `os.Stdin`.

### `NewMultiSource(sources ...Source) Source`

Combines multiple sources into a single source. Lines from all sources are merged into one channel. Useful for tailing several log files simultaneously.

## Usage

```go
src, err := source.NewFileSource("/var/log/app.log")
if err != nil {
    log.Fatal(err)
}

for line := range src.Lines() {
    fmt.Println(line)
}
```

To read from multiple files at once:

```go
src1, _ := source.NewFileSource("/var/log/app.log")
src2, _ := source.NewFileSource("/var/log/app-error.log")

for line := range source.NewMultiSource(src1, src2).Lines() {
    fmt.Println(line)
}
```

All sources close their output channel when the underlying reader reaches EOF.
