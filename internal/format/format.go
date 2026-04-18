// Package format provides field value formatting rules for log entries.
package format

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines how to format a specific field.
type Rule struct {
	Field  string
	Format string // e.g. "%.2f", "%05d", "%s"
}

// Formatter applies format rules to log entries.
type Formatter struct {
	rules []Rule
}

// New returns a Formatter with the given rules.
func New(rules []Rule) *Formatter {
	return &Formatter{rules: rules}
}

// Apply returns a new entry with formatted field values.
// Fields that do not exist or cannot be formatted are left unchanged.
func (f *Formatter) Apply(e parser.Entry) parser.Entry {
	if len(f.rules) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		out[k] = v
	}
	for _, r := range f.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		formatted := fmt.Sprintf(r.Format, v)
		out[r.Field] = formatted
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: out, Raw: e.Raw}
}

// ParseRules parses specs of the form "field=format", e.g. "latency=%.3f".
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		idx := strings.Index(s, "=")
		if idx < 0 {
			return nil, fmt.Errorf("format: invalid spec %q: expected field=format", s)
		}
		field := strings.TrimSpace(s[:idx])
		fmt_ := strings.TrimSpace(s[idx+1:])
		if field == "" {
			return nil, fmt.Errorf("format: empty field name in spec %q", s)
		}
		if fmt_ == "" {
			return nil, fmt.Errorf("format: empty format string in spec %q", s)
		}
		rules = append(rules, Rule{Field: field, Format: fmt_})
	}
	return rules, nil
}
