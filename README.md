# logslice

Stream and filter structured JSON logs from multiple sources with a unified query syntax.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

---

## Usage

Pipe logs directly or point logslice at a file, socket, or remote source and filter using the query syntax:

```bash
# Filter logs from a file where level is "error"
logslice --source ./app.log 'level == "error"'

# Stream from multiple sources and match on a field
logslice --source ./service-a.log --source ./service-b.log 'status >= 500'

# Tail a running process and pretty-print matches
logslice --tail --source ./app.log --pretty 'user_id == "abc123" && latency > 200'
```

### Query Syntax

| Operator | Example |
|----------|---------|
| `==`     | `level == "warn"` |
| `!=`     | `env != "prod"` |
| `>` / `<` | `latency > 300` |
| `&&` / `\|\|` | `level == "error" && service == "auth"` |

---

## Features

- Reads from files, stdin, or multiple sources simultaneously
- Filters structured JSON log lines using a simple expression syntax
- Optional pretty-printed output for human-readable inspection
- Minimal dependencies, single binary

---

## License

MIT © [yourusername](https://github.com/yourusername)