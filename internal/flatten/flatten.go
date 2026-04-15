// Package flatten provides utilities for flattening nested JSON log entries
// into a single-level map using dot-separated key paths.
package flatten

import (
	"fmt"
	"strings"
)

// Flattener collapses nested map fields into dot-notation keys.
type Flattener struct {
	separator string
	maxDepth  int
}

// Option configures a Flattener.
type Option func(*Flattener)

// WithSeparator sets the key separator (default ".").
func WithSeparator(sep string) Option {
	return func(f *Flattener) {
		if sep != "" {
			f.separator = sep
		}
	}
}

// WithMaxDepth limits how deep the flattening recurses (0 = unlimited).
func WithMaxDepth(d int) Option {
	return func(f *Flattener) {
		if d >= 0 {
			f.maxDepth = d
		}
	}
}

// New creates a Flattener with the given options.
func New(opts ...Option) *Flattener {
	f := &Flattener{separator: "."}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Apply flattens a nested map[string]any entry in place, returning a new map.
func (f *Flattener) Apply(entry map[string]any) map[string]any {
	out := make(map[string]any, len(entry))
	f.flatten("", entry, out, 0)
	return out
}

func (f *Flattener) flatten(prefix string, src map[string]any, dst map[string]any, depth int) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = prefix + f.separator + k
		}
		if nested, ok := v.(map[string]any); ok && (f.maxDepth == 0 || depth < f.maxDepth) {
			f.flatten(key, nested, dst, depth+1)
		} else {
			dst[key] = v
		}
	}
}

// ParseKeys returns a sorted list of all dot-notation keys present in the entry.
func ParseKeys(entry map[string]any) []string {
	keys := make([]string, 0, len(entry))
	collectKeys("", entry, &keys)
	return keys
}

func collectKeys(prefix string, m map[string]any, out *[]string) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, k)
		}
		if nested, ok := v.(map[string]any); ok {
			collectKeys(key, nested, out)
		} else {
			*out = append(*out, key)
		}
	}
}

// HasPrefix reports whether any key in the flat entry starts with the given prefix.
func HasPrefix(entry map[string]any, prefix string) bool {
	for k := range entry {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}
