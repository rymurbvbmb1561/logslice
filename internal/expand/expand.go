package expand

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Expander extracts a nested JSON string field and merges its keys into the
// parent entry. This is useful when a log field contains an escaped JSON
// payload that should be promoted to top-level fields.
type Expander struct {
	fields    []string
	prefix    string
	overwrite bool
}

// Option configures an Expander.
type Option func(*Expander)

// WithPrefix sets a prefix applied to every key extracted from the nested
// JSON value.
func WithPrefix(p string) Option {
	return func(e *Expander) { e.prefix = p }
}

// WithOverwrite controls whether expanded keys overwrite existing keys in the
// parent entry. Default is false (existing keys are preserved).
func WithOverwrite(v bool) Option {
	return func(e *Expander) { e.overwrite = v }
}

// New creates an Expander that will expand the named fields.
func New(fields []string, opts ...Option) *Expander {
	e := &Expander{fields: fields}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Apply expands each configured field in entry, returning a new map.
// If a field does not exist or is not a valid JSON object string it is left
// untouched and no error is returned.
func (e *Expander) Apply(entry map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{}, len(entry))
	for k, v := range entry {
		out[k] = v
	}

	for _, field := range e.fields {
		raw, ok := out[field]
		if !ok {
			continue
		}
		s, ok := raw.(string)
		if !ok {
			continue
		}
		s = strings.TrimSpace(s)
		if !strings.HasPrefix(s, "{") {
			continue
		}
		var nested map[string]interface{}
		if err := json.Unmarshal([]byte(s), &nested); err != nil {
			// not valid JSON — leave field as-is
			continue
		}
		delete(out, field)
		for k, v := range nested {
			key := fmt.Sprintf("%s%s", e.prefix, k)
			if _, exists := out[key]; exists && !e.overwrite {
				continue
			}
			out[key] = v
		}
	}
	return out, nil
}
