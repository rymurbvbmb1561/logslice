// Package filter provides time-range and field-value filtering for log entries.
package filter

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Filter defines criteria for matching log entries.
type Filter struct {
	// From is the inclusive start of the time range. Nil means no lower bound.
	From *time.Time
	// To is the inclusive end of the time range. Nil means no upper bound.
	To *time.Time
	// Fields is a map of field names to required values.
	Fields map[string]string
}

// Match reports whether the given entry satisfies all criteria in f.
func Match(e parser.Entry, f Filter) bool {
	if f.From != nil && e.Timestamp.Before(*f.From) {
		return false
	}
	if f.To != nil && e.Timestamp.After(*f.To) {
		return false
	}
	for k, v := range f.Fields {
		got, ok := e.Fields[k]
		if !ok {
			return false
		}
		gotStr, ok := got.(string)
		if !ok {
			return false
		}
		if gotStr != v {
			return false
		}
	}
	return true
}

// Apply returns only the entries from entries that satisfy f.
func Apply(entries []parser.Entry, f Filter) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if Match(e, f) {
			result = append(result, e)
		}
	}
	return result
}
