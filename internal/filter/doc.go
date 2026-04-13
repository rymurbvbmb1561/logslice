// Package filter implements log entry filtering for logslice.
//
// It provides the Options type for specifying filter criteria and the
// Match and Apply functions for evaluating log entries against those criteria.
//
// Supported filter criteria:
//
//   - Time range: filter entries by a start time (TimeFrom) and/or end time (TimeTo).
//   - Field equality: filter entries where a named field equals an expected string value.
//
// Example usage:
//
//	opts := filter.Options{
//	    TimeFrom: start,
//	    TimeTo:   end,
//	    Fields:   map[string]string{"level": "error"},
//	}
//	matched := filter.Apply(entries, opts)
package filter
