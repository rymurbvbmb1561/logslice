// Package format provides a Formatter that applies Go fmt-style format strings
// to named fields in a log entry.
//
// Rules are specified as "field=format" pairs, e.g.:
//
//	"latency=%.3f"
//	"status_code=%d"
//
// Fields that are absent or whose type is incompatible with the format verb
// are left unchanged. The Formatter never mutates the original entry.
package format
