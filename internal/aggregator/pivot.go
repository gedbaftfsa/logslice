package aggregator

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// Pivot groups entries by a key field and aggregates a value field using a
// summary function (count, sum, avg, min, max).
type Pivot struct {
	mu      sync.Mutex
	keyField string
	valField string
	op       string
	buckets  map[string]*pivotBucket
}

type pivotBucket struct {
	count int
	sum   float64
	min   float64
	max   float64
}

// NewPivot creates a Pivot aggregator. op must be one of: count, sum, avg, min, max.
func NewPivot(keyField, valField, op string) (*Pivot, error) {
	switch op {
	case "count", "sum", "avg", "min", "max":
	default:
		return nil, fmt.Errorf("pivot: unknown op %q; want count|sum|avg|min|max", op)
	}
	return &Pivot{
		keyField: keyField,
		valField:  valField,
		op:        op,
		buckets:   make(map[string]*pivotBucket),
	}, nil
}

// Record ingests a raw JSON log line.
func (p *Pivot) Record(raw []byte) {
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return
	}
	keyRaw, ok := m[p.keyField]
	if !ok {
		return
	}
	key := fmt.Sprintf("%v", keyRaw)

	var val float64
	if p.op != "count" {
		v, ok := m[p.valField]
		if !ok {
			return
		}
		switch n := v.(type) {
		case float64:
			val = n
		default:
			return
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	b, exists := p.buckets[key]
	if !exists {
		b = &pivotBucket{min: val, max: val}
		p.buckets[key] = b
	}
	b.count++
	b.sum += val
	if val < b.min {
		b.min = val
	}
	if val > b.max {
		b.max = val
	}
}

// Snapshot returns a sorted slice of JSON objects representing the pivot table.
func (p *Pivot) Snapshot() []map[string]interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	keys := make([]string, 0, len(p.buckets))
	for k := range p.buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make([]map[string]interface{}, 0, len(keys))
	for _, k := range keys {
		b := p.buckets[k]
		var agg float64
		switch p.op {
		case "count":
			agg = float64(b.count)
		case "sum":
			agg = b.sum
		case "avg":
			if b.count > 0 {
				agg = b.sum / float64(b.count)
			}
		case "min":
			agg = b.min
		case "max":
			agg = b.max
		}
		result = append(result, map[string]interface{}{
			p.keyField: k,
			p.op:       agg,
		})
	}
	return result
}

// Reset clears all accumulated data.
func (p *Pivot) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.buckets = make(map[string]*pivotBucket)
}
