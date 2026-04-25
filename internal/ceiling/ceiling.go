package ceiling

import (
	"fmt"
	"math"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a ceiling (round-up) operation on a numeric field to the nearest multiple.
type Rule struct {
	Field    string
	Multiple float64
}

// Applier rounds numeric field values up to the nearest multiple.
type Applier struct {
	rules []Rule
}

// New returns an Applier with the given rules.
func New(rules []Rule) *Applier {
	return &Applier{rules: rules}
}

// Apply rounds each matching field value up to the nearest multiple.
// Non-numeric or missing fields are left unchanged.
func (a *Applier) Apply(entry parser.Entry) parser.Entry {
	if len(a.rules) == 0 {
		return entry
	}
	out := entry.Clone()
	for _, r := range a.rules {
		v, ok := out.Fields[r.Field]
		if !ok {
			continue
		}
		f, ok := toFloat(v)
		if !ok {
			continue
		}
		multiple := r.Multiple
		if multiple <= 0 {
			multiple = 1
		}
		out.Fields[r.Field] = math.Ceil(f/multiple) * multiple
	}
	return out
}

func toFloat(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case string:
		var f float64
		_, err := fmt.Sscanf(val, "%g", &f)
		return f, err == nil
	}
	return 0, false
}

// ParseRules parses specs of the form "field=multiple", e.g. "latency=100".
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("ceiling: invalid spec %q: expected field=multiple", s)
		}
		field := strings.TrimSpace(parts[0])
		if field == "" {
			return nil, fmt.Errorf("ceiling: empty field name in spec %q", s)
		}
		var multiple float64
		if _, err := fmt.Sscanf(strings.TrimSpace(parts[1]), "%g", &multiple); err != nil {
			return nil, fmt.Errorf("ceiling: invalid multiple in spec %q: %w", s, err)
		}
		if multiple <= 0 {
			return nil, fmt.Errorf("ceiling: multiple must be positive in spec %q", s)
		}
		rules = append(rules, Rule{Field: field, Multiple: multiple})
	}
	return rules, nil
}
