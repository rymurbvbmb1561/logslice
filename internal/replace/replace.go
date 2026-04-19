package replace

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule describes a single find-and-replace operation on a log entry field.
type Rule struct {
	Field       string
	Pattern     *regexp.Regexp
	Replacement string
}

// Replacer applies replacement rules to log entries.
type Replacer struct {
	rules []Rule
}

// New returns a Replacer with the given rules.
func New(rules []Rule) *Replacer {
	return &Replacer{rules: rules}
}

// Apply returns a new entry with replacement rules applied.
// Non-string fields and missing fields are left unchanged.
func (r *Replacer) Apply(e parser.Entry) parser.Entry {
	if len(r.rules) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		out[k] = v
	}
	for _, rule := range r.rules {
		v, ok := out[rule.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[rule.Field] = rule.Pattern.ReplaceAllString(s, rule.Replacement)
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: out, Raw: e.Raw}
}

// ParseRules parses specs of the form "field/pattern/replacement".
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "/", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("replace: invalid spec %q: expected field/pattern/replacement", spec)
		}
		field, pattern, replacement := parts[0], parts[1], parts[2]
		if field == "" {
			return nil, fmt.Errorf("replace: empty field in spec %q", spec)
		}
		if pattern == "" {
			return nil, fmt.Errorf("replace: empty pattern in spec %q", spec)
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("replace: invalid pattern %q: %w", pattern, err)
		}
		rules = append(rules, Rule{Field: field, Pattern: re, Replacement: replacement})
	}
	return rules, nil
}
