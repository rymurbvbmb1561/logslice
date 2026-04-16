// Package typecast provides field type coercion for log entries.
package typecast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule describes a single field-to-type casting rule.
type Rule struct {
	Field string
	Type  string // "int", "float", "bool", "string"
}

// Caster applies type casting rules to log entries.
type Caster struct {
	rules []Rule
}

// New returns a Caster with the given rules.
func New(rules []Rule) *Caster {
	return &Caster{rules: rules}
}

// Apply returns a new entry with fields cast according to the rules.
// Fields that cannot be cast are left unchanged.
func (c *Caster) Apply(e parser.Entry) parser.Entry {
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
		if cast, err := castValue(v, r.Type); err == nil {
			out[r.Field] = cast
		}
	}
	return parser.Entry{Fields: out, Raw: e.Raw}
}

func castValue(v any, typ string) (any, error) {
	s := fmt.Sprintf("%v", v)
	switch strings.ToLower(typ) {
	case "int":
		return strconv.ParseInt(s, 10, 64)
	case "float":
		return strconv.ParseFloat(s, 64)
	case "bool":
		return strconv.ParseBool(s)
	case "string":
		return s, nil
	default:
		return nil, fmt.Errorf("unknown type %q", typ)
	}
}

// ParseRules parses specs like "field=type" into Rules.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid typecast rule %q: expected field=type", s)
		}
		rules = append(rules, Rule{Field: parts[0], Type: parts[1]})
	}
	return rules, nil
}
