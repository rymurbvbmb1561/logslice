// Package filter provides functionality for filtering log entries
// based on time ranges and field value conditions.
package filter

import (
	"time"
)

// Options holds the filtering criteria for log entries.
type Options struct {
	// TimeFrom filters entries at or after this time. Zero value means no lower bound.
	TimeFrom time.Time
	// TimeTo filters entries at or before this time. Zero value means no upper bound.
	TimeTo time.Time
	// Fields is a map of field name to expected value for equality filtering.
	Fields map[string]string
}

// Entry represents a parsed log entry with a timestamp and arbitrary fields.
type Entry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// Match reports whether the entry satisfies all conditions in opts.
func Match(entry Entry, opts Options) bool {
	if !opts.TimeFrom.IsZero() && entry.Timestamp.Before(opts.TimeFrom) {
		return false
	}
	if !opts.TimeTo.IsZero() && entry.Timestamp.After(opts.TimeTo) {
		return false
	}
	for key, expected := range opts.Fields {
		val, ok := entry.Fields[key]
		if !ok {
			return false
		}
		strVal, ok := val.(string)
		if !ok {
			return false
		}
		if strVal != expected {
			return false
		}
	}
	return true
}

// Apply filters a slice of entries according to opts and returns matching entries.
func Apply(entries []Entry, opts Options) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if Match(e, opts) {
			result = append(result, e)
		}
	}
	return result
}
