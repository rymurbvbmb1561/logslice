// Package cli implements the command-line interface for logslice.
//
// It is responsible for parsing flags and arguments, wiring together
// the reader, parser, filter, and output packages into a single
// processing pipeline.
//
// Supported flags:
//
//	--from    Start of the time range (RFC3339 or YYYY-MM-DD)
//	--to      End of the time range (RFC3339 or YYYY-MM-DD)
//	--fields  Comma-separated key=value pairs to filter log fields
//	--format  Output format: raw (default), json, or text
//
// A positional argument may be supplied as the input file path.
// If omitted, logslice reads from standard input.
package cli
