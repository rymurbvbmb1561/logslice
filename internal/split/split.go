package split

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Splitter splits a single log entry field containing a delimited string
// into multiple named fields.
type Splitter struct {
	rules []Rule
}

// Rule describes how to split one field into many.
type Rule struct {
	Source  string
	Targets []string
	Delim   string
}

// New returns a Splitter configured with the given rules.
func New(rules []Rule) *Splitter {
	return &Splitter{rules: rules}
}

// Apply splits fields in entry according to configured rules.
// The original source field is preserved. A new entry is returned.
func (s *Splitter) Apply(entry parser.Entry) parser.Entry {
	if len(s.rules) == 0 {
		return entry
	}
	out := entry.Clone()
	for _, r := range s.rules {
		v, ok := out.Fields[r.Source]
		if !ok {
			continue
		}
		str, ok := v.(string)
		if !ok {
			continue
		}
		parts := strings.SplitN(str, r.Delim, len(r.Targets))
		for i, target := range r.Targets {
			if i < len(parts) {
				out.Fields[target] = parts[i]
			}
		}
	}
	return out
}

// ParseRules parses specs of the form "source:target1,target2|delim".
// Example: "msg:level,text| "
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		rule, err := parseRule(spec)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// parseRule parses a single spec of the form "source:target1,target2|delim".
func parseRule(spec string) (Rule, error) {
	pipe := strings.LastIndex(spec, "|")
	if pipe < 0 {
		return Rule{}, fmt.Errorf("split: missing '|' delimiter in spec %q", spec)
	}
	delim := spec[pipe+1:]
	if delim == "" {
		return Rule{}, fmt.Errorf("split: empty delimiter in spec %q", spec)
	}
	left := spec[:pipe]
	colon := strings.Index(left, ":")
	if colon < 0 {
		return Rule{}, fmt.Errorf("split: missing ':' in spec %q", spec)
	}
	source := left[:colon]
	if source == "" {
		return Rule{}, fmt.Errorf("split: empty source field in spec %q", spec)
	}
	targets := strings.Split(left[colon+1:], ",")
	if len(targets) == 0 || (len(targets) == 1 && targets[0] == "") {
		return Rule{}, fmt.Errorf("split: no target fields in spec %q", spec)
	}
	return Rule{Source: source, Targets: targets, Delim: delim}, nil
}
