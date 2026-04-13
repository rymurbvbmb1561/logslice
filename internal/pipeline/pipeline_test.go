package pipeline_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/stats"
)

func makeJSONLine(ts, level, msg string) string {
	b, _ := json.Marshal(map[string]string{"time": ts, "level": level, "msg": msg})
	return string(b)
}

func TestRun_FiltersByTime(t *testing.T) {
	lines := strings.Join([]string{
		makeJSONLine("2024-01-01T10:00:00Z", "info", "early"),
		makeJSONLine("2024-01-01T12:00:00Z", "info", "match"),
		makeJSONLine("2024-01-01T14:00:00Z", "info", "late"),
	}, "\n")

	r := reader.NewStdinReader(strings.NewReader(lines))
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)
	s := stats.New()

	from := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)

	err := pipeline.Run(pipeline.Config{
		Reader: r,
		Filter: filter.Options{From: &from, To: &to},
		Writer: w,
		Stats:  s,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Matched() != 1 {
		t.Errorf("expected 1 matched, got %d", s.Matched())
	}
	if s.Read() != 3 {
		t.Errorf("expected 3 read, got %d", s.Read())
	}
	if !strings.Contains(buf.String(), "match") {
		t.Errorf("expected output to contain 'match', got: %s", buf.String())
	}
}

func TestRun_ParseErrorsSkipped(t *testing.T) {
	lines := strings.Join([]string{
		"not json at all",
		makeJSONLine("2024-01-01T12:00:00Z", "info", "ok"),
	}, "\n")

	r := reader.NewStdinReader(strings.NewReader(lines))
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)
	s := stats.New()

	err := pipeline.Run(pipeline.Config{
		Reader: r,
		Filter: filter.Options{},
		Writer: w,
		Stats:  s,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ParseErrors() != 1 {
		t.Errorf("expected 1 parse error, got %d", s.ParseErrors())
	}
	if s.Matched() != 1 {
		t.Errorf("expected 1 matched, got %d", s.Matched())
	}
}
