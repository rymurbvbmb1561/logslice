package timebucket_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/timebucket"
)

func makeEntry(ts time.Time) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Fields:    map[string]any{"msg": "test"},
		Raw:       `{"msg":"test"}`,
	}
}

func TestNew_DefaultInterval(t *testing.T) {
	b := timebucket.New()
	if b == nil {
		t.Fatal("expected non-nil Bucketer")
	}
}

func TestRecord_ZeroTimestampIgnored(t *testing.T) {
	b := timebucket.New()
	b.Record(parser.Entry{Fields: map[string]any{"msg": "no ts"}})
	if len(b.Buckets()) != 0 {
		t.Errorf("expected 0 buckets, got %d", len(b.Buckets()))
	}
}

func TestRecord_GroupsByInterval(t *testing.T) {
	b := timebucket.New(timebucket.WithInterval(time.Minute))
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	b.Record(makeEntry(base))
	b.Record(makeEntry(base.Add(30 * time.Second)))
	b.Record(makeEntry(base.Add(90 * time.Second)))

	buckets := b.Buckets()
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}

	key1 := base.Truncate(time.Minute)
	if buckets[key1] != 2 {
		t.Errorf("expected 2 in first bucket, got %d", buckets[key1])
	}

	key2 := base.Add(90 * time.Second).Truncate(time.Minute)
	if buckets[key2] != 1 {
		t.Errorf("expected 1 in second bucket, got %d", buckets[key2])
	}
}

func TestWithInterval_ZeroIgnored(t *testing.T) {
	b := timebucket.New(timebucket.WithInterval(0))
	base := time.Date(2024, 1, 1, 12, 0, 30, 0, time.UTC)
	b.Record(makeEntry(base))
	// default 1m should still apply
	key := base.Truncate(time.Minute)
	if b.Buckets()[key] != 1 {
		t.Errorf("expected default interval to apply")
	}
}

func TestWriteSummary_ContainsBucketTime(t *testing.T) {
	b := timebucket.New(timebucket.WithInterval(time.Hour))
	ts := time.Date(2024, 6, 15, 9, 45, 0, 0, time.UTC)
	b.Record(makeEntry(ts))

	var buf bytes.Buffer
	b.WriteSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "2024-06-15T09:00:00Z") {
		t.Errorf("expected truncated hour in output, got:\n%s", out)
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected count 1 in output, got:\n%s", out)
	}
}
