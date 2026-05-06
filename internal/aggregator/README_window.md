# aggregator/window

Provides a **sliding time window** over structured JSON log entries.

## Overview

`Window` buffers `json.RawMessage` entries and automatically evicts those that
fall outside a configurable duration relative to the most-recently added entry.

## Usage

```go
w := aggregator.NewWindow(30 * time.Second)

// Add entries as they arrive from a source.
w.Add(rawJSON, "time") // "time" is the RFC3339 timestamp field name

// Retrieve all entries currently inside the window.
entries := w.Entries()
```

## Timestamp resolution

The field named by the second argument to `Add` is parsed as an **RFC3339**
string. If the field is absent, empty, or unparseable, `time.Now()` is used as
a fallback so the entry is always accepted.

## Eviction

Eviction happens lazily — on every call to `Add`, `Entries`, or `Len`. Entries
whose timestamp is older than `newest - duration` are dropped.

## Types

| Symbol | Description |
|--------|-------------|
| `NewWindow(d time.Duration) *Window` | Create a new window with the given duration |
| `(*Window).Add(raw, field)` | Insert an entry, evict stale ones |
| `(*Window).Entries() []json.RawMessage` | Return all live entries |
| `(*Window).Len() int` | Count of live entries |
