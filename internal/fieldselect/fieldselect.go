package fieldselect

import "github.com/user/logslice/internal/parser"

// Selector keeps or drops fields from log entries.
type Selector struct {
	fields  map[string]struct{}
	inverted bool // if true, drop listed fields instead of keeping
}

// Option configures a Selector.
type Option func(*Selector)

// WithFields sets the fields to keep.
func WithFields(fields []string) Option {
	return func(s *Selector) {
		for _, f := range fields {
			s.fields[f] = struct{}{}
		}
	}
}

// WithDrop inverts the selector so listed fields are dropped.
func WithDrop() Option {
	return func(s *Selector) {
		s.inverted = true
	}
}

// New creates a Selector with the given options.
func New(opts ...Option) *Selector {
	s := &Selector{fields: make(map[string]struct{})}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Apply returns a new entry containing only the selected fields.
func (s *Selector) Apply(e parser.Entry) parser.Entry {
	if len(s.fields) == 0 {
		return e
	}
	out := parser.Entry{
		Raw:    e.Raw,
		Fields: make(map[string]any, len(e.Fields)),
	}
	for k, v := range e.Fields {
		_, listed := s.fields[k]
		if (!s.inverted && listed) || (s.inverted && !listed) {
			out.Fields[k] = v
		}
	}
	return out
}

// Fields returns the set of field names configured on the selector.
func (s *Selector) Fields() []string {
	out := make([]string, 0, len(s.fields))
	for f := range s.fields {
		out = append(out, f)
	}
	return out
}

// ParseFields splits a comma-separated field list.
func ParseFields(spec string) []string {
	if spec == "" {
		return nil
	}
	var out []string
	start := 0
	for i := 0; i <= len(spec); i++ {
		if i == len(spec) || spec[i] == ',' {
			f := spec[start:i]
			if f != "" {
				out = append(out, f)
			}
			start = i + 1
		}
	}
	return out
}
