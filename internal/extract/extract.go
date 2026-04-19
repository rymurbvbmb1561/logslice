package extract

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a regex extraction rule: apply pattern to source field,
// write named capture groups as new fields.
type Rule struct {
	Source  string
	Pattern *regexp.Regexp
}

// Extractor applies extraction rules to log entries.
type Extractor struct {
	rules []Rule
}

// New returns an Extractor with the given rules.
func New(rules []Rule) *Extractor {
	return &Extractor{rules: rules}
}

// Apply runs all extraction rules against the entry and returns a new entry
// with any captured fields merged in. The original entry is not mutated.
func (e *Extractor) Apply(entry parser.Entry) parser.Entry {
	if len(e.rules) == 0 {
		return entry
	}
	out := entry.Clone()
	for _, r := range e.rules {
		val, ok := out.Fields[r.Source]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		match := r.Pattern.FindStringSubmatch(s)
		if match == nil {
			continue
		}
		for i, name := range r.Pattern.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			out.Fields[name] = match[i]
		}
	}
	return out
}

// ParseRules parses specs of the form "field=pattern" where pattern must
// contain at least one named capture group (?P<name>...).
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		idx := strings.IndexByte(spec, '=')
		if idx < 1 {
			return nil, fmt.Errorf("extract: invalid spec %q: expected field=pattern", spec)
		}
		field := spec[:idx]
		patStr := spec[idx+1:]
		if patStr == "" {
			return nil, fmt.Errorf("extract: empty pattern for field %q", field)
		}
		pat, err := regexp.Compile(patStr)
		if err != nil {
			return nil, fmt.Errorf("extract: invalid pattern for field %q: %w", field, err)
		}
		rules = append(rules, Rule{Source: field, Pattern: pat})
	}
	return rules, nil
}
