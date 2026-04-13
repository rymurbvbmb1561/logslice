package pipeline_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/stats"
)

func TestRun_NilReader(t *testing.T) {
	err := pipeline.Run(pipeline.Config{
		Writer: output.NewWriter(strings.NewReader(""), output.FormatJSON),
		Stats:  stats.New(),
	})
	if err == nil || err.Error() != pipeline.ErrNoReader.Error() {
		t.Errorf("expected ErrNoReader, got %v", err)
	}
}

func TestRun_NilWriter(t *testing.T) {
	err := pipeline.Run(pipeline.Config{
		Reader: reader.NewStdinReader(strings.NewReader("")),
		Stats:  stats.New(),
	})
	if err == nil || err.Error() != pipeline.ErrNoWriter.Error() {
		t.Errorf("expected ErrNoWriter, got %v", err)
	}
}

func TestRun_NilStats(t *testing.T) {
	var buf strings.Builder
	err := pipeline.Run(pipeline.Config{
		Reader: reader.NewStdinReader(strings.NewReader("")),
		Writer: output.NewWriter(&buf, output.FormatJSON),
		Filter: filter.Options{},
	})
	if err == nil || err.Error() != pipeline.ErrNoStats.Error() {
		t.Errorf("expected ErrNoStats, got %v", err)
	}
}
