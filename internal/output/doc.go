// Package output provides utilities for writing filtered log entries
// to an io.Writer in multiple formats.
//
// Supported formats:
//
//   - FormatJSON  – re-serialises the parsed fields map as a JSON object.
//   - FormatText  – writes a human-readable line with an ISO-8601 timestamp
//     followed by the fields map.
//   - FormatRaw   – writes the original raw log line unchanged, which is
//     useful when the input should pass through unmodified.
//
// Example usage:
//
//	w := output.NewStdoutWriter(output.FormatJSON)
//	for _, entry := range filtered {
//		w.Write(entry)
//	}
package output
