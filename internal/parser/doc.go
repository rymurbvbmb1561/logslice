// Package parser provides functionality for parsing structured (JSON) log lines
// into LogEntry values that carry a resolved timestamp and a map of raw fields.
//
// Usage:
//
//	entry, err := parser.ParseLine(`{"time":"2024-01-02T15:04:05Z","level":"info","msg":"ok"}`)
//	if err != nil {
//	    // handle parse error (empty line, invalid JSON, etc.)
//	}
//	fmt.Println(entry.Timestamp, entry.Fields["msg"])
//
// Timestamp detection
//
// The parser recognises the following field names as timestamp sources,
// checked in order: "time", "timestamp", "ts", "@timestamp".
// Both string values (RFC 3339 and common variants) and numeric Unix epoch
// values (float64 as decoded by encoding/json) are supported.
package parser
