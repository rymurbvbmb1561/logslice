package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry represents a single parsed log line with its raw fields.
type LogEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// TimeFormats lists common timestamp formats to attempt during parsing.
var TimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05.000Z07:00",
}

// ParseLine parses a single JSON log line into a LogEntry.
// It attempts to extract a timestamp from common field names.
func ParseLine(line string) (*LogEntry, error) {
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	fields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	entry := &LogEntry{
		Fields: fields,
		Raw:    line,
	}

	for _, key := range []string{"time", "timestamp", "ts", "@timestamp"} {
		if val, ok := fields[key]; ok {
			if ts, err := parseTimestamp(val); err == nil {
				entry.Timestamp = ts
				break
			}
		}
	}

	return entry, nil
}

// parseTimestamp attempts to parse a value as a time.Time using known formats.
func parseTimestamp(val interface{}) (time.Time, error) {
	switch v := val.(type) {
	case string:
		for _, format := range TimeFormats {
			if t, err := time.Parse(format, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("unrecognized timestamp format: %s", v)
	case float64:
		return time.Unix(int64(v), 0).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported timestamp type: %T", val)
	}
}
