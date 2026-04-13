package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

var baseTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func makeEntry(ts time.Time, fields map[string]interface{}) filter.Entry {
	return filter.Entry{Timestamp: ts, Fields: fields, Raw: "{}"}
}

func TestMatch_TimeFrom(t *testing.T) {
	entry := makeEntry(baseTime, nil)
	opts := filter.Options{TimeFrom: baseTime.Add(-time.Hour)}
	if !filter.Match(entry, opts) {
		t.Error("expected entry to match TimeFrom filter")
	}
	opts.TimeFrom = baseTime.Add(time.Hour)
	if filter.Match(entry, opts) {
		t.Error("expected entry to not match TimeFrom filter")
	}
}

func TestMatch_TimeTo(t *testing.T) {
	entry := makeEntry(baseTime, nil)
	opts := filter.Options{TimeTo: baseTime.Add(time.Hour)}
	if !filter.Match(entry, opts) {
		t.Error("expected entry to match TimeTo filter")
	}
	opts.TimeTo = baseTime.Add(-time.Hour)
	if filter.Match(entry, opts) {
		t.Error("expected entry to not match TimeTo filter")
	}
}

func TestMatch_FieldFilter(t *testing.T) {
	entry := makeEntry(baseTime, map[string]interface{}{"level": "error", "service": "api"})
	opts := filter.Options{Fields: map[string]string{"level": "error"}}
	if !filter.Match(entry, opts) {
		t.Error("expected entry to match field filter")
	}
	opts.Fields["level"] = "info"
	if filter.Match(entry, opts) {
		t.Error("expected entry to not match field filter")
	}
}

func TestMatch_MissingField(t *testing.T) {
	entry := makeEntry(baseTime, map[string]interface{}{"level": "error"})
	opts := filter.Options{Fields: map[string]string{"service": "api"}}
	if filter.Match(entry, opts) {
		t.Error("expected entry to not match when field is missing")
	}
}

func TestApply(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(baseTime.Add(-2*time.Hour), map[string]interface{}{"level": "info"}),
		makeEntry(baseTime, map[string]interface{}{"level": "error"}),
		makeEntry(baseTime.Add(2*time.Hour), map[string]interface{}{"level": "error"}),
	}
	opts := filter.Options{
		TimeFrom: baseTime.Add(-time.Hour),
		TimeTo:   baseTime.Add(time.Hour),
		Fields:   map[string]string{"level": "error"},
	}
	result := filter.Apply(entries, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Timestamp != baseTime {
		t.Error("unexpected entry in result")
	}
}
