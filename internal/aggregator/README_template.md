# Template Aggregator

The `Template` aggregator renders a [Go template](https://pkg.go.dev/text/template)
against each JSON log entry and stores the result as a new string field.

## Usage

```go
tmpl, err := aggregator.NewTemplate("summary", "[{{.level}}] {{.msg}}")
if err != nil {
    log.Fatal(err)
}

for _, line := range lines {
    tmpl.Record(line)
}

for _, result := range tmpl.Results() {
    fmt.Println(string(result))
}
```

## Constructor

```
NewTemplate(field, templateStr string) (*Template, error)
```

| Argument      | Description                                      |
|---------------|--------------------------------------------------|
| `field`       | Output field name to store the rendered string   |
| `templateStr` | Go template string evaluated against each entry  |

Returns an error if `field` or `templateStr` is empty, or if the template
fails to parse.

## Behaviour

- Each JSON log entry is decoded into a `map[string]any`.
- The template is executed with the map as its data context.
- The rendered string is stored under `field` in the output object.
- Missing keys default to `<no value>` (Go's zero-value template behaviour).
- Invalid JSON lines return an error and are not stored.

## Example

Input:
```json
{"level":"error","msg":"disk full","host":"web-01"}
```

Template: `[{{.level}}] {{.msg}} on {{.host}}`

Output:
```json
{"level":"error","msg":"disk full","host":"web-01","summary":"[error] disk full on web-01"}
```
