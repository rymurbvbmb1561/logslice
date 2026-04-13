package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/stats"
)

func TestNew_InitialCountsAreZero(t *testing.T) {
	c := stats.New()
	if c.LinesRead != 0 || c.LinesParsed != 0 || c.LinesMatched != 0 || c.ParseErrors != 0 {
		t.Error("expected all counters to be zero on creation")
	}
}

func TestRecordRead(t *testing.T) {
	c := stats.New()
	c.RecordRead()
	c.RecordRead()
	if c.LinesRead != 2 {
		t.Errorf("expected LinesRead=2, got %d", c.LinesRead)
	}
}

func TestRecordParsed(t *testing.T) {
	c := stats.New()
	c.RecordParsed()
	if c.LinesParsed != 1 {
		t.Errorf("expected LinesParsed=1, got %d", c.LinesParsed)
	}
}

func TestRecordMatched(t *testing.T) {
	c := stats.New()
	c.RecordMatched()
	c.RecordMatched()
	c.RecordMatched()
	if c.LinesMatched != 3 {
		t.Errorf("expected LinesMatched=3, got %d", c.LinesMatched)
	}
}

func TestRecordParseError(t *testing.T) {
	c := stats.New()
	c.RecordParseError()
	if c.ParseErrors != 1 {
		t.Errorf("expected ParseErrors=1, got %d", c.ParseErrors)
	}
}

func TestElapsed_IsPositive(t *testing.T) {
	c := stats.New()
	if c.Elapsed() < 0 {
		t.Error("elapsed duration should be non-negative")
	}
}

func TestPrint_ContainsExpectedFields(t *testing.T) {
	c := stats.New()
	for i := 0; i < 10; i++ {
		c.RecordRead()
	}
	for i := 0; i < 8; i++ {
		c.RecordParsed()
	}
	for i := 0; i < 5; i++ {
		c.RecordMatched()
	}
	for i := 0; i < 2; i++ {
		c.RecordParseError()
	}

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	for _, want := range []string{"10", "8", "5", "2", "elapsed"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q; got:\n%s", want, out)
		}
	}
}

func TestPrint_NoParseErrorsLine_WhenZero(t *testing.T) {
	c := stats.New()
	c.RecordRead()
	c.RecordParsed()
	c.RecordMatched()

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	if strings.Contains(out, "parse errors") {
		t.Errorf("expected no 'parse errors' line when count is zero; got:\n%s", out)
	}
}
