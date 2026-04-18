// Package prefix provides a transformer that prepends a static string
// to a named string field in a log entry.
package prefix

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule describes a single prefix operation.
type Rule struct {
	Field  string
	Prefix string
}

// Prefixer applies prefix rules to log entries.
type Prefixer struct {
	rules []Rule
}

// New returns a Prefixer that applies the given rules.
func New(rules []Rule) *Prefixer {
	return &Prefixer{rules: rules}
}

// Apply returns a new entry with prefix rules applied.
// Non-string fields and missing fields are left unchanged.
func (p *Prefixer) Apply(e parser.Entry) parser.Entry {
	if len(p.rules) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		out[k] = v
	}
	for _, r := range p.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[r.Field] = r.Prefix + s
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: out, Raw: e.Raw}
}

// ParseRules parses specs of the form "field=prefix".
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("prefix: invalid spec %q: expected field=prefix", s)
		}
		if parts[0] == "" {
			return nil, fmt.Errorf("prefix: empty field name in spec %q", s)
		}
		rules = append(rules, Rule{Field: parts[0], Prefix: parts[1]})
	}
	return rules, nil
}
