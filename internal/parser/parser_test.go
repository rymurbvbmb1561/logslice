package parser

import (
	"testing"
	"time"
)

func TestParseLine_ValidJSON(t *testing.T) {
	line := `{"time":"2024-03-15T10:30:00Z","level":"info","msg":"server started"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", entry.Fields["level"])
	}
	expected := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entry.Timestamp)
	}
}

func TestParseLine_UnixTimestamp(t *testing.T) {
	line := `{"ts":1710498600,"msg":"hello"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp from unix ts field")
	}
}

func TestParseLine_InvalidJSON(t *testing.T) {
	_, err := ParseLine("not json at all")
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	_, err := ParseLine("")
	if err == nil {
		t.Error("expected error for empty line, got nil")
	}
}

func TestParseLine_NoTimestampField(t *testing.T) {
	line := `{"level":"warn","msg":"missing timestamp"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !entry.Timestamp.IsZero() {
		t.Errorf("expected zero timestamp, got %v", entry.Timestamp)
	}
}

func TestParseLine_AtTimestampKey(t *testing.T) {
	line := `{"@timestamp":"2024-06-01T08:00:00Z","service":"api"}`
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp.Year() != 2024 {
		t.Errorf("expected year 2024, got %d", entry.Timestamp.Year())
	}
}
