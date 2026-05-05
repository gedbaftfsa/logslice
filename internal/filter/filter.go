package filter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Expr represents a single filter expression (e.g. "level=error", "status>=400")
type Expr struct {
	Field    string
	Operator string
	Value    string
}

// Filter holds a slice of expressions that are ANDed together.
type Filter struct {
	Exprs []Expr
}

var operators = []string{"!=", ">=", "<=", ">", "<", "="}

// Parse parses a query string like "level=error status>=400" into a Filter.
func Parse(query string) (*Filter, error) {
	if query == "" {
		return &Filter{}, nil
	}
	parts := strings.Fields(query)
	exprs := make([]Expr, 0, len(parts))
	for _, part := range parts {
		expr, err := parseExpr(part)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}
	return &Filter{Exprs: exprs}, nil
}

func parseExpr(s string) (Expr, error) {
	for _, op := range operators {
		if idx := strings.Index(s, op); idx > 0 {
			return Expr{
				Field:    s[:idx],
				Operator: op,
				Value:    s[idx+len(op):],
			}, nil
		}
	}
	return Expr{}, fmt.Errorf("invalid filter expression: %q", s)
}

// Match returns true if the JSON log line satisfies all filter expressions.
func (f *Filter) Match(line []byte) bool {
	if len(f.Exprs) == 0 {
		return true
	}
	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return false
	}
	for _, expr := range f.Exprs {
		if !matchExpr(record, expr) {
			return false
		}
	}
	return true
}

func matchExpr(record map[string]interface{}, expr Expr) bool {
	val, ok := record[expr.Field]
	if !ok {
		return false
	}
	actual := fmt.Sprintf("%v", val)
	switch expr.Operator {
	case "=":
		return actual == expr.Value
	case "!=":
		return actual != expr.Value
	}
	// Numeric comparisons
	actualNum, err1 := strconv.ParseFloat(actual, 64)
	wantNum, err2 := strconv.ParseFloat(expr.Value, 64)
	if err1 != nil || err2 != nil {
		return false
	}
	switch expr.Operator {
	case ">":
		return actualNum > wantNum
	case ">=":
		return actualNum >= wantNum
	case "<":
		return actualNum < wantNum
	case "<=":
		return actualNum <= wantNum
	}
	return false
}
