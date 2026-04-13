// Package pipeline provides the top-level processing pipeline for logslice.
//
// It connects the reader, parser, filter, output, and stats packages into a
// single, testable unit. Callers construct a [Config] with the desired
// components and invoke [Run] to process log lines end-to-end.
//
// Typical usage:
//
//	err := pipeline.Run(pipeline.Config{
//		Reader: r,
//		Filter: filter.Options{From: &from, To: &to},
//		Writer: w,
//		Stats:  s,
//	})
//
// Run returns the first fatal error encountered (e.g. I/O failure). Parse
// errors on individual lines are counted in Stats and skipped rather than
// aborting the run.
package pipeline
