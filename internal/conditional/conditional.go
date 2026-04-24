package conditional

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule describes a conditional field assignment: if field matches value, set target to result.
type Rule struct {
	Field  string
	Match  string
	Target string
	Value  string
}

// Applier applies conditional rules to log entries.
type Applier struct {
	rules []Rule
}

// New returns an Applier with the given rules.
func New(rules []Rule) *Applier {
	return &Applier{rules: rules}
}

// Apply evaluates each rule against the entry and returns a (possibly modified) copy.
// If the entry's Field equals the rule's Match value, Target is set to Value.
// The original entry is never mutated.
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
		if fmt.Sprintf("%v", v) == r.Match {
			out.Fields[r.Target] = r.Value
		}
	}
	return out
}

// ParseRules parses specs of the form "field=match:target=value".
// Multiple specs may be provided as separate slice elements.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		if spec == "" {
			continue
		}
		// Expected format: field=match:target=value
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("conditional: invalid spec %q: expected 'field=match:target=value'", spec)
		}
		lhs := strings.SplitN(parts[0], "=", 2)
		if len(lhs) != 2 || lhs[0] == "" {
			return nil, fmt.Errorf("conditional: invalid condition in spec %q", spec)
		}
		rhs := strings.SplitN(parts[1], "=", 2)
		if len(rhs) != 2 || rhs[0] == "" {
			return nil, fmt.Errorf("conditional: invalid assignment in spec %q", spec)
		}
		rules = append(rules, Rule{
			Field:  lhs[0],
			Match:  lhs[1],
			Target: rhs[0],
			Value:  rhs[1],
		})
	}
	return rules, nil
}
