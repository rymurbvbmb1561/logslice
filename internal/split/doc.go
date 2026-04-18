// Package split provides a log entry transformer that splits a single
// delimited string field into multiple named target fields.
//
// Rules are specified as "source:target1,target2|delimiter" strings.
// For example, "addr:host,port|:" splits the "addr" field on ":" into
// "host" and "port". The original source field is preserved in the output.
//
// Usage:
//
//	rules, err := split.ParseRules(specs)
//	s := split.New(rules)
//	out := s.Apply(entry)
package split
