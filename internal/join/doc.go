// Package join provides a log entry transformer that concatenates multiple
// fields into a single target field using a configurable separator.
//
// Rules are specified as strings in the form:
//
//	"target=source1,source2[|separator]"
//
// If the separator is omitted, a single space is used. Missing source fields
// are silently skipped. The original entry is never mutated.
package join
