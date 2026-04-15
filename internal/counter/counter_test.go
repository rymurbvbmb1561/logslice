package counter_test

import (
	"testing"

	"github.com/logslice/logslice/internal/counter"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestNew_InitialCountsEmpty(t *testing.T) {
	c := counter.New("level")
	if got := c.Counts(); len(got) != 0 {
		t.Fatalf("expected empty counts, got %v", got)
	}
}

func TestRecord_CountsByFieldValue(t *testing.T) {
	c := counter.New("level")
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "error"}))

	counts := c.Counts()
	if len(counts) != 2 {
		t.Fatalf("expected 2 distinct values, got %d", len(counts))
	}
	if counts[0].Value != "info" || counts[0].Count != 2 {
		t.Errorf("unexpected first entry: %+v", counts[0])
	}
	if counts[1].Value != "error" || counts[1].Count != 1 {
		t.Errorf("unexpected second entry: %+v", counts[1])
	}
}

func TestRecord_MissingFieldIgnored(t *testing.T) {
	c := counter.New("level")
	c.Record(makeEntry(map[string]any{"msg": "hello"}))
	if got := c.Counts(); len(got) != 0 {
		t.Fatalf("expected empty counts, got %v", got)
	}
}

func TestTop_ReturnsSortedByCount(t *testing.T) {
	c := counter.New("level")
	for i := 0; i < 3; i++ {
		c.Record(makeEntry(map[string]any{"level": "debug"}))
	}
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "error"}))

	top := c.Top(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0].Value != "debug" {
		t.Errorf("expected debug first, got %s", top[0].Value)
	}
	if top[1].Value != "info" {
		t.Errorf("expected info second, got %s", top[1].Value)
	}
}

func TestWithLimit_CapsDistinctValues(t *testing.T) {
	c := counter.New("level", counter.WithLimit(2))
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "error"}))
	c.Record(makeEntry(map[string]any{"level": "debug"}))

	if got := len(c.Counts()); got != 2 {
		t.Errorf("expected 2 distinct values due to limit, got %d", got)
	}
}

func TestTop_NZero_ReturnsAll(t *testing.T) {
	c := counter.New("status")
	c.Record(makeEntry(map[string]any{"status": "200"}))
	c.Record(makeEntry(map[string]any{"status": "404"}))
	c.Record(makeEntry(map[string]any{"status": "500"}))

	if got := len(c.Top(0)); got != 3 {
		t.Errorf("expected all 3, got %d", got)
	}
}
