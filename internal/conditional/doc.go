// Package conditional provides conditional field assignment for log entries.
//
// A Rule specifies: if a given field equals a match value, set a target field
// to a result value. Rules are evaluated in order against each entry.
//
// Example spec format (for ParseRules):
//
//	"level=error:alert=true"
//
// This sets the "alert" field to "true" whenever the "level" field equals "error".
//
// Entries are never mutated; Apply always returns a cloned copy.
package conditional
