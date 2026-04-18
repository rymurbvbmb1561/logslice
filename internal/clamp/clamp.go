// Package clamp provides a processor that clamps numeric field values
// to a specified [min, max] range.
package clamp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a clamping rule for a single field.
type Rule struct {
	Field string
	Min   float64
	Max   float64
}

// Clamper applies clamping rules to log entries.
type Clamper struct {
	rules []Rule
}

// New returns a new Clamper with the given rules.
func New(rules []Rule) *Clamper {
	return &Clamper{rules: rules}
}

// Apply returns a copy of the entry with numeric fields clamped per rules.
func (c *Clamper) Apply(e parser.Entry) parser.Entry {
	if len(c.rules) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		out[k] = v
	}
	for _, r := range c.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		f, err := toFloat(v)
		if err != nil {
			continue
		}
		if f < r.Min {
			f = r.Min
		} else if f > r.Max {
			f = r.Max
		}
		out[r.Field] = f
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: out, Raw: e.Raw}
}

func toFloat(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	}
	return 0, fmt.Errorf("not numeric")
}

// ParseRules parses specs like "field=min:max" into Rules.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("clamp: invalid spec %q: expected field=min:max", s)
		}
		bounds := strings.SplitN(parts[1], ":", 2)
		if len(bounds) != 2 {
			return nil, fmt.Errorf("clamp: invalid spec %q: expected min:max", s)
		}
		min, err := strconv.ParseFloat(bounds[0], 64)
		if err != nil {
			return nil, fmt.Errorf("clamp: invalid min in %q: %w", s, err)
		}
		max, err := strconv.ParseFloat(bounds[1], 64)
		if err != nil {
			return nil, fmt.Errorf("clamp: invalid max in %q: %w", s, err)
		}
		if min > max {
			return nil, fmt.Errorf("clamp: min > max in %q", s)
		}
		rules = append(rules, Rule{Field: parts[0], Min: min, Max: max})
	}
	return rules, nil
}
