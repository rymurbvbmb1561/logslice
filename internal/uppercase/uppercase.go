// Package uppercase provides a transformer that converts string field values to uppercase.
package uppercase

import (
	"fmt"
	"strings"

	"github.com/your-org/logslice/internal/parser"
)

// Transformer converts specified string fields to uppercase.
type Transformer struct {
	fields []string
}

// Option configures a Transformer.
type Option func(*Transformer)

// WithFields sets the fields to uppercase.
func WithFields(fields []string) Option {
	return func(t *Transformer) {
		t.fields = fields
	}
}

// New creates a new uppercase Transformer.
func New(opts ...Option) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply returns a copy of the entry with the configured fields uppercased.
func (t *Transformer) Apply(e parser.Entry) parser.Entry {
	if len(t.fields) == 0 {
		return e
	}
	out := make(parser.Entry, len(e))
	for k, v := range e {
		out[k] = v
	}
	for _, f := range t.fields {
		v, ok := out[f]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[f] = strings.ToUpper(s)
	}
	return out
}

// ParseFields parses a comma-separated list of field names.
func ParseFields(spec string) ([]string, error) {
	if strings.TrimSpace(spec) == "" {
		return nil, nil
	}
	parts := strings.Split(spec, ",")
	var fields []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("uppercase: empty field name in spec %q", spec)
		}
		fields = append(fields, p)
	}
	return fields, nil
}
