// Package template provides a simple field-interpolation engine for
// formatting structured log entries as human-readable strings.
//
// Templates use {fieldname} placeholders which are replaced with the
// corresponding field values from a parsed log entry. Missing fields
// are rendered as "<nil>".
//
// Example:
//
//	applier, err := template.ParseTemplate("{time} [{level}] {msg}")
//	if err != nil {
//		log.Fatal(err)
//	}
//	out, err := applier.Apply(entry)
package template
