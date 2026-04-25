// Package ceiling provides an Applier that rounds numeric log entry field
// values up to the nearest multiple of a configured step.
//
// This is useful for bucketing continuous metrics (e.g. latency, byte counts)
// into discrete intervals for grouping or display.
//
// Rules are specified as "field=multiple" strings, for example:
//
//	"latency_ms=100"  // rounds latency_ms up to nearest 100
//	"response_size=512"  // rounds response_size up to nearest 512
//
// Fields that are missing, non-numeric, or already on a boundary are left
// unchanged. The original entry is never mutated.
package ceiling
