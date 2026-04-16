// Package compute provides field computation for log entries.
//
// It supports deriving new numeric fields from existing ones using
// simple arithmetic expressions. Rules are specified as strings of
// the form:
//
//	target=source<op>operand
//
// Supported operators: *, /, +, -
//
// Example:
//
//	duration_ms=duration*1000
//	rate=count/interval
//
// Fields that are missing or non-numeric are silently skipped.
// The original entry is never mutated.
package compute
