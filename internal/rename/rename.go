package rename

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule maps an old field name to a new field name.
type Rule struct {
	From string
	To   string
}

// Renamer renames fields in log entries.
type Renamer struct {
	rules []Rule
}

// New returns a Renamer applying the given rules.
func New(rules []Rule) *Renamer {
	return &Renamer{rules: rules}
}

// Apply returns a new entry with fields renamed according to the rules.
// If a source field does not exist the rule is silently skipped.
func (r *Renamer) Apply(entry parser.Entry) parser.Entry {
	if len(r.rules) == 0 {
		return entry
	}
	out := make(parser.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, rule := range r.rules {
		val, ok := out[rule.From]
		if !ok {
			continue
		}
		out[rule.To] = val
		delete(out, rule.From)
	}
	return out
}

// ParseRules parses specs of the form "old=new" into a slice of Rule.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("rename: invalid spec %q: expected old=new", s)
		}
		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		if from == "" {
			return nil, fmt.Errorf("rename: empty source field in spec %q", s)
		}
		if to == "" {
			return nil, fmt.Errorf("rename: empty target field in spec %q", s)
		}
		rules = append(rules, Rule{From: from, To: to})
	}
	return rules, nil
}
