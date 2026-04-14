package redact

import (
	"regexp"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Redactor replaces sensitive field values with a placeholder string.
type Redactor struct {
	fields      map[string]struct{}
	patterns    []*regexp.Regexp
	placeholder string
}

// Option configures a Redactor.
type Option func(*Redactor)

// WithFields marks specific field names whose values should be redacted.
func WithFields(fields ...string) Option {
	return func(r *Redactor) {
		for _, f := range fields {
			r.fields[f] = struct{}{}
		}
	}
}

// WithPatterns adds regex patterns; any field value matching a pattern is redacted.
func WithPatterns(patterns ...*regexp.Regexp) Option {
	return func(r *Redactor) {
		r.patterns = append(r.patterns, patterns...)
	}
}

// WithPlaceholder sets the replacement string (default: "[REDACTED]").
func WithPlaceholder(p string) Option {
	return func(r *Redactor) {
		r.placeholder = p
	}
}

// New creates a Redactor with the given options.
func New(opts ...Option) *Redactor {
	r := &Redactor{
		fields:      make(map[string]struct{}),
		placeholder: "[REDACTED]",
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Apply returns a copy of entry with sensitive values replaced.
func (r *Redactor) Apply(entry parser.Entry) parser.Entry {
	out := make(parser.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for k := range r.fields {
		if _, ok := out[k]; ok {
			out[k] = r.placeholder
		}
	}
	for k, v := range out {
		str, ok := v.(string)
		if !ok {
			continue
		}
		for _, pat := range r.patterns {
			if pat.MatchString(str) {
				out[k] = pat.ReplaceAllString(str, r.placeholder)
				break
			}
		}
	}
	return out
}

// ParseFields parses a comma-separated list of field names.
func ParseFields(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
