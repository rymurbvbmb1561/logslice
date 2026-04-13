// Package cli provides the command-line interface for logslice.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

// Config holds all CLI flag values.
type Config struct {
	From       string
	To         string
	Fields     []string
	Format     string
	InputFile  string
}

// Run parses arguments and executes the main pipeline.
func Run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}
	return execute(cfg)
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	from := fs.String("from", "", "start time (RFC3339 or Unix timestamp)")
	to := fs.String("to", "", "end time (RFC3339 or Unix timestamp)")
	fields := fs.String("fields", "", "comma-separated key=value field filters")
	format := fs.String("format", "raw", fmt.Sprintf("output format (%s)", strings.Join(output.KnownFormats, ", ")))

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		From:      *from,
		To:        *to,
		Format:    *format,
		InputFile: fs.Arg(0),
	}

	if *fields != "" {
		cfg.Fields = strings.Split(*fields, ",")
	}

	return cfg, nil
}

func execute(cfg *Config) error {
	var r reader.LineReader
	var err error

	if cfg.InputFile != "" {
		r, err = reader.NewFileReader(cfg.InputFile)
		if err != nil {
			return fmt.Errorf("opening input: %w", err)
		}
	} else {
		r = reader.NewStdinReader()
	}

	fmt := output.ParseFormat(cfg.Format)
	w := output.NewStdoutWriter(fmt)

	var f filter.Filter
	if cfg.From != "" {
		t, err := parseTime(cfg.From)
		if err != nil {
			return fmt.Errorf("invalid --from: %w", err)
		}
		f.From = &t
	}
	if cfg.To != "" {
		t, err := parseTime(cfg.To)
		if err != nil {
			return fmt.Errorf("invalid --to: %w", err)
		}
		f.To = &t
	}
	f.Fields = parseFieldFilters(cfg.Fields)

	lines, err := r.Lines()
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	var entries []parser.Entry
	for _, line := range lines {
		e, err := parser.ParseLine(line)
		if err != nil {
			continue
		}
		entries = append(entries, e)
	}

	matched := filter.Apply(entries, f)
	return w.WriteAll(matched)
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("unsupported time format; use RFC3339 or YYYY-MM-DD")
}

func parseFieldFilters(pairs []string) map[string]string {
	if len(pairs) == 0 {
		return nil
	}
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m
}
