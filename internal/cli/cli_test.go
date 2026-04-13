package cli

import (
	"testing"
	"time"
)

func TestParseTime_RFC3339(t *testing.T) {
	got, err := parseTime("2024-01-15T10:30:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTime_DateOnly(t *testing.T) {
	got, err := parseTime("2024-03-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Year() != 2024 || got.Month() != 3 || got.Day() != 1 {
		t.Errorf("unexpected date: %v", got)
	}
}

func TestParseTime_Invalid(t *testing.T) {
	_, err := parseTime("not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid time string")
	}
}

func TestParseFieldFilters_Valid(t *testing.T) {
	result := parseFieldFilters([]string{"level=error", "service=api"})
	if result["level"] != "error" {
		t.Errorf("expected level=error, got %q", result["level"])
	}
	if result["service"] != "api" {
		t.Errorf("expected service=api, got %q", result["service"])
	}
}

func TestParseFieldFilters_Empty(t *testing.T) {
	result := parseFieldFilters(nil)
	if result != nil {
		t.Errorf("expected nil map for empty input")
	}
}

func TestParseFieldFilters_MalformedEntry(t *testing.T) {
	result := parseFieldFilters([]string{"level=error", "badentry"})
	if _, ok := result["badentry"]; ok {
		t.Error("malformed entry should not be added to filter map")
	}
	if result["level"] != "error" {
		t.Errorf("valid entry should still be parsed")
	}
}

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "raw" {
		t.Errorf("expected default format 'raw', got %q", cfg.Format)
	}
	if cfg.From != "" || cfg.To != "" {
		t.Error("expected empty from/to by default")
	}
}

func TestParseFlags_AllOptions(t *testing.T) {
	cfg, err := parseFlags([]string{
		"--from", "2024-01-01T00:00:00Z",
		"--to", "2024-01-02T00:00:00Z",
		"--fields", "level=error,service=api",
		"--format", "json",
		"myfile.log",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.From != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected From: %q", cfg.From)
	}
	if cfg.Format != "json" {
		t.Errorf("unexpected Format: %q", cfg.Format)
	}
	if cfg.InputFile != "myfile.log" {
		t.Errorf("unexpected InputFile: %q", cfg.InputFile)
	}
	if len(cfg.Fields) != 2 {
		t.Errorf("expected 2 field filters, got %d", len(cfg.Fields))
	}
}
