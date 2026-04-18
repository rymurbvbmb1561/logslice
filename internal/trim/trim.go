// Package trim provides a processor that trims leading and/or trailing
// whitespace (or a custom cutset) from string fields in a log entry.
package trim

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Direction controls which end(s) of the string are trimmed.
type Direction int

const (
	Both  Direction = iota
	Left            // leading characters only
	Right           // trailing characters only
)

// Rule describes a single trim operation.
type Rule struct {
	Field  string
	Cutset string // empty means whitespace
	Dir    Direction
}

// Trimmer applies trim rules to log entries.
type Trimmer struct {
	rules []Rule
}

// New returns a Trimmer configured with the given rules.
func New(rules []Rule) *Trimmer {
	return &Trimmer{rules: rules}
}

// Apply returns a new entry with string fields trimmed according to the rules.
// The original entry is never mutated.
func (t *Trimmer) Apply(e parser.Entry) parser.Entry {
	if len(t.rules) == 0 {
		return e
	}
	out := parser.Entry{Timestamp: e.Timestamp, Raw: e.Raw, Fields: make(map[string]any, len(e.Fields))}
	for k, v := range e.Fields {
		out.Fields[k] = v
	}
	for _, r := range t.rules {
		v, ok := out.Fields[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out.Fields[r.Field] = applyTrim(s, r)
	}
	return out
}

func applyTrim(s string, r Rule) string {
	if r.Cutset == "" {
		switch r.Dir {
		case Left:
			return strings.TrimLeft(s, " \t\n\r")
		case Right:
			return strings.TrimRight(s, " \t\n\r")
		default:
			return strings.TrimSpace(s)
		}
	}
	switch r.Dir {
	case Left:
		return strings.TrimLeft(s, r.Cutset)
	case Right:
		return strings.TrimRight(s, r.Cutset)
	default:
		return strings.Trim(s, r.Cutset)
	}
}
