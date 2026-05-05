# filter

Package `filter` provides query-string parsing and JSON log line matching for logslice.

## Query Syntax

A query is a space-separated list of filter expressions. All expressions are **ANDed** together.

### Operators

| Operator | Meaning           |
|----------|-------------------|
| `=`      | Equal             |
| `!=`     | Not equal         |
| `>`      | Greater than      |
| `>=`     | Greater or equal  |
| `<`      | Less than         |
| `<=`     | Less or equal     |

### Examples

```
level=error
level=error status>=500
status!=200 latency>100
```

## Usage

```go
f, err := filter.Parse("level=error status>=400")
if err != nil {
    log.Fatal(err)
}

for _, line := range logLines {
    if f.Match(line) {
        fmt.Println(string(line))
    }
}
```

Numeric comparisons (`>`, `>=`, `<`, `<=`) require the field value in the JSON to be parseable as a float64.
