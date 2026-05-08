# Compute

The `Compute` aggregator applies an arithmetic operation to a numeric field in
each log entry and emits the result as a new field.

## Usage

```go
c, err := aggregator.NewCompute(srcField, dstField, op, operand)
if err != nil {
    log.Fatal(err)
}
for _, entry := range entries {
    c.Record(entry)
}
results := c.Results()
```

## Supported Operations

| Op      | Description                          |
|---------|--------------------------------------|
| `add`   | `dst = src + operand`                |
| `sub`   | `dst = src - operand`                |
| `mul`   | `dst = src * operand`                |
| `div`   | `dst = src / operand` (operand ≠ 0)  |
| `abs`   | `dst = |src|` (operand ignored)      |
| `round` | `dst = round(src)` (operand ignored) |

## Notes

- Entries where `srcField` is absent or non-numeric are silently skipped.
- The original entry is not mutated; a new JSON object is emitted.
- Call `Reset()` to clear accumulated results between passes.
