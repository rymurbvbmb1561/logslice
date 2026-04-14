// Package truncator provides field-level value truncation for log entries
// before they are written to output. Long string values can be capped at a
// configurable maximum byte length, with an optional ellipsis suffix.
package truncator

import "github.com/user/logslice/internal/parser"

const defaultEllipsis = "..."

// Options controls truncation behaviour.
type Options struct {
	// MaxLen is the maximum number of bytes allowed for any string field value.
	// Values longer than MaxLen are truncated. A MaxLen of 0 disables truncation.
	MaxLen int
	// Ellipsis is appended to truncated values. Defaults to "..." when empty.
	Ellipsis string
}

// Truncator applies field-value truncation to log entries.
type Truncator struct {
	opts Options
}

// New returns a Truncator configured with opts.
func New(opts Options) *Truncator {
	if opts.Ellipsis == "" {
		opts.Ellipsis = defaultEllipsis
	}
	return &Truncator{opts: opts}
}

// Apply returns a copy of entry with all string field values truncated to
// MaxLen bytes. If MaxLen is 0 the entry is returned unchanged.
func (t *Truncator) Apply(entry parser.Entry) parser.Entry {
	if t.opts.MaxLen <= 0 {
		return entry
	}

	truncated := make(map[string]interface{}, len(entry.Fields))
	for k, v := range entry.Fields {
		s, ok := v.(string)
		if !ok {
			truncated[k] = v
			continue
		}
		if len(s) > t.opts.MaxLen {
			cutoff := t.opts.MaxLen
			// Ensure we don't cut in the middle of a multi-byte rune.
			for cutoff > 0 && !isRuneBoundary(s, cutoff) {
				cutoff--
			}
			s = s[:cutoff] + t.opts.Ellipsis
		}
		truncated[k] = s
	}

	return parser.Entry{
		Timestamp: entry.Timestamp,
		Raw:       entry.Raw,
		Fields:    truncated,
	}
}

// isRuneBoundary reports whether index i is on a UTF-8 rune boundary in s.
func isRuneBoundary(s string, i int) bool {
	if i == 0 || i == len(s) {
		return true
	}
	return (s[i] & 0xC0) != 0x80
}
