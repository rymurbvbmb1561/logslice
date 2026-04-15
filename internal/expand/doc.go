// Package expand provides an Expander that promotes nested JSON string fields
// into the top-level log entry map.
//
// Some log pipelines store a secondary JSON payload as an escaped string inside
// a field (e.g. "payload": "{\"user\":\"alice\"}"). Expander detects such
// fields, parses the inner JSON object, and merges its keys directly into the
// parent entry — optionally applying a prefix to avoid key collisions.
//
// Example usage:
//
//	e := expand.New(
//		[]string{"payload"},
//		expand.WithPrefix("payload_"),
//		expand.WithOverwrite(false),
//	)
//	out, err := e.Apply(entry)
package expand
