package compute

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a computed field: target = expr (e.g. "duration_ms=response_time*1000").
type Rule struct {
	Target string
	Source string
	Op     string
	Operand float64
}

// Applier applies compute rules to log entries.
type Applier struct {
	rules []Rule
}

// New returns an Applier with the given rules.
func New(rules []Rule) *Applier {
	return &Applier{rules: rules}
}

// Apply returns a new entry with computed fields added.
func (a *Applier) Apply(e parser.Entry) parser.Entry {
	if len(a.rules) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields)+len(a.rules))
	for k, v := range e.Fields {
		out[k] = v
	}
	for _, r := range a.rules {
		v, ok := toFloat(e.Fields[r.Source])
		if !ok {
			continue
		}
		var result float64
		switch r.Op {
		case "*":
			result = v * r.Operand
		case "/":
			if r.Operand == 0 {
				continue
			}
			result = v / r.Operand
		case "+":
			result = v + r.Operand
		case "-":
			result = v - r.Operand
		default:
			continue
		}
		out[r.Target] = result
	}
	return parser.Entry{Fields: out, Raw: e.Raw, Timestamp: e.Timestamp}
}

func toFloat(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	}
	return 0, false
}

// ParseRules parses specs like "duration_ms=response_time*1000".
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	ops := []string{"*", "/", "+", "-"}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("compute: invalid spec %q: expected target=source<op>operand", spec)
		}
		target := strings.TrimSpace(parts[0])
		expr := strings.TrimSpace(parts[1])
		var matched bool
		for _, op := range ops {
			idx := strings.LastIndex(expr, op)
			if idx <= 0 {
				continue
			}
			source := strings.TrimSpace(expr[:idx])
			operandStr := strings.TrimSpace(expr[idx+1:])
			operand, err := strconv.ParseFloat(operandStr, 64)
			if err != nil {
				continue
			}
			rules = append(rules, Rule{Target: target, Source: source, Op: op, Operand: operand})
			matched = true
			break
		}
		if !matched {
			return nil, fmt.Errorf("compute: invalid expr %q in spec %q", expr, spec)
		}
	}
	return rules, nil
}
