# cli

The `cli` package implements the command-line interface for `logslice`.

## Usage

```
logslice [options] [file...]
```

If no files are provided, `logslice` reads from **stdin**.

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `-filter` | `""` | Filter expression, e.g. `level=error` |
| `-format` | `json` | Output format: `json`, `pretty`, `text` |
| `-fields` | `""` | Comma-separated fields to include in output |

## Examples

```bash
# Filter error logs from a file
logslice -filter 'level=error' app.log

# Read from stdin, output pretty-printed
cat app.log | logslice -format pretty

# Select specific fields from multiple files
logslice -fields 'ts,level,msg' -filter 'level=warn' a.log b.log
```
