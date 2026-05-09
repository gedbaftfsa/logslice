# aggregator/redact

The `Redact` processor masks sensitive fields in JSON log lines by replacing their values with a configurable placeholder string.

## Usage

```go
r, err := aggregator.NewRedact([]string{"password", "token"}, "[REDACTED]")
if err != nil {
    log.Fatal(err)
}

out, err := r.Record(line)
```

## Constructor

```go
func NewRedact(fields []string, placeholder string) (*Redact, error)
```

- `fields` — one or more JSON field names to redact (required).
- `placeholder` — the string to substitute; defaults to `"***"` if empty.

Returns an error if `fields` is empty.

## Methods

| Method | Description |
|--------|-------------|
| `Record(line []byte) ([]byte, error)` | Redacts configured fields in the JSON line. Passes through unchanged if the field is absent or the line is not valid JSON. |
| `Snapshot() []byte` | Returns a JSON summary with `redacted_fields`, `lines_redacted`, and `placeholder`. |
| `Reset()` | Clears the redaction counter. |

## Behaviour

- If a line is not valid JSON it is passed through unchanged and no error is returned.
- If none of the target fields are present in a line the original bytes are returned unmodified and the counter is not incremented.
- Fields not listed in `fields` are left untouched.
