# tail

The `tail` package provides live file-following functionality, similar to
`tail -f` on Unix systems.

## Overview

`Follow` opens a file and streams each new line to a channel as the file
grows. It uses polling rather than OS-specific inotify/kqueue APIs to keep
the implementation portable.

## Usage

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

lines, err := tail.Follow(ctx, "/var/log/app.log", tail.Options{
    PollInterval: 100 * time.Millisecond,
})
if err != nil {
    log.Fatal(err)
}

for line := range lines {
    fmt.Println(line)
}
```

## Options

| Field          | Default | Description                              |
|----------------|---------|------------------------------------------|
| `PollInterval` | 200 ms  | How often to re-check the file for data  |

## Notes

- The channel is closed when the context is cancelled or a read error occurs.
- Blank lines are silently skipped.
- Existing content is emitted before waiting for new writes.
