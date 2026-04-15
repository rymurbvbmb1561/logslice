package aggregate_test

import {
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/aggregate"
	"github.com/user/logslice/internal/parser"
}

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       "{}",
		Fields:    fields,
	}
}

func TestNew_InitialCountsEmpty(t *testing.T) {
	c := aggregate.New("level")
	if got := c.Total(); got != 0 {
		t.Fatalf("expected total 0, got %d", got)
	}
	if len(c.Counts()) != 0 {
		t.Fatal("expected empty counts map")
	}
}

func TestRecord_CountsByFieldValue(t *testing.T) {
	c := aggregate.New("level")
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "error"}))

	counts := c.Counts()
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", counts["error"])
	}
	if c.Total() != 3 {
		t.Errorf("expected total=3, got %d", c.Total())
	}
}

func TestRecord_MissingFieldCountedSeparately(t *testing.T) {
	c := aggregate.New("level")
	c.Record(makeEntry(map[string]any{"msg": "hello"}))

	counts := c.Counts()
	if counts["(missing)"] != 1 {
		t.Errorf("expected (missing)=1, got %d", counts["(missing)"])
	}
}

func TestWriteSummary_OutputContainsValues(t *testing.T) {
	c := aggregate.New("level")
	c.Record(makeEntry(map[string]any{"level": "warn"}))
	c.Record(makeEntry(map[string]any{"level": "info"}))
	c.Record(makeEntry(map[string]any{"level": "info"}))

	var buf bytes.Buffer
	c.WriteSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "info") {
		t.Error("expected output to contain 'info'")
	}
	if !strings.Contains(out, "warn") {
		t.Error("expected output to contain 'warn'")
	}
	// info (count=2) should appear before warn (count=1)
	infoIdx := strings.Index(out, "info")
	warnIdx := strings.Index(out, "warn")
	if infoIdx > warnIdx {
		t.Error("expected info to appear before warn (higher count first)")
	}
}

func TestCounts_ReturnsCopy(t *testing.T) {
	c := aggregate.New("level")
	c.Record(makeEntry(map[string]any{"level": "info"}))

	counts := c.Counts()
	counts["info"] = 999

	if c.Counts()["info"] != 1 {
		t.Error("Counts() should return a copy, not a reference")
	}
}
