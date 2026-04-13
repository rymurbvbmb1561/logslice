// Package pipeline wires together the reader, parser, filter, output,
// and stats components into a single processing pipeline.
package pipeline

import (
	"fmt"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/stats"
)

// Config holds all parameters needed to run the pipeline.
type Config struct {
	Reader  *reader.Reader
	Filter  filter.Options
	Writer  *output.Writer
	Stats   *stats.Stats
}

// Run reads every line from cfg.Reader, parses it, applies the filter,
// and writes matching entries via cfg.Writer. Progress is recorded in
// cfg.Stats. The first write error terminates the run.
func Run(cfg Config) error {
	lines, err := cfg.Reader.Lines()
	if err != nil {
		return fmt.Errorf("pipeline: reading lines: %w", err)
	}

	var entries []parser.Entry
	for _, line := range lines {
		cfg.Stats.RecordRead()

		entry, err := parser.ParseLine(line)
		if err != nil {
			cfg.Stats.RecordParseError()
			continue
		}
		cfg.Stats.RecordParsed()

		if filter.Match(entry, cfg.Filter) {
			cfg.Stats.RecordMatched()
			entries = append(entries, entry)
		}
	}

	if err := cfg.Writer.WriteAll(entries); err != nil {
		return fmt.Errorf("pipeline: writing output: %w", err)
	}
	return nil
}
