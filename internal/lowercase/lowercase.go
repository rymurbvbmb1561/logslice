package lowercase

import (
	"fmt"
	"strings"

	"github.com/your-org/logslice/internal/parser"
)

// Applier lowercases string field values.
type Applier struct {
	fields map[string]struct{}
}

// WithFields returns an Applier that lowercases only the given fields.
func WithFields(fields []string) *Applier {
	m := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		m[f] = struct{}{}
	}
	return &Applier{fields: m}
}

// New returns an Applier with the provided fields.
func New(fields []string) *Applier {
	return WithFields(fields)
}

// Apply returns a new entry with the specified string fields lowercased.
func (a *Applier) Apply(e parser.Entry) parser.Entry {
	if len(a.fields) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		out[k] = v
	}
	for f := range a.fields {
		v, ok := out[f]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[f] = strings.ToLower(s)
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: out, Raw: e.Raw}
}

// ParseFields parses a comma-separated list of field names.
func ParseFields(spec string) ([]string, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, nil
	}
	parts := strings.Split(spec, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("lowercase: empty field name in spec %q", spec)
		}
		out = append(out, p)
	}
	return out, nil
}
