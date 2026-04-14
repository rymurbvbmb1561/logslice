// Package dedup provides log entry deduplication based on a configurable field or full-line hashing.
package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Entry represents a parsed log entry passed through the pipeline.
type Entry map[string]interface{}

// Deduplicator tracks seen entries and reports whether a new entry is a duplicate.
type Deduplicator struct {
	seen  map[string]struct{}
	field string // if empty, hash the full raw line
}

// New creates a Deduplicator. If field is non-empty, deduplication is based
// on the value of that field; otherwise the entire JSON representation is hashed.
func New(field string) *Deduplicator {
	return &Deduplicator{
		seen:  make(map[string]struct{}),
		field: field,
	}
}

// IsDuplicate returns true if the entry has been seen before, and records it
// if it has not. It is safe to call sequentially (not concurrency-safe).
func (d *Deduplicator) IsDuplicate(entry Entry) (bool, error) {
	key, err := d.keyFor(entry)
	if err != nil {
		return false, err
	}
	if _, exists := d.seen[key]; exists {
		return true, nil
	}
	d.seen[key] = struct{}{}
	return false, nil
}

// Seen returns the number of unique entries recorded so far.
func (d *Deduplicator) Seen() int {
	return len(d.seen)
}

// Reset clears the deduplication state.
func (d *Deduplicator) Reset() {
	d.seen = make(map[string]struct{})
}

func (d *Deduplicator) keyFor(entry Entry) (string, error) {
	if d.field != "" {
		val, ok := entry[d.field]
		if !ok {
			return "", fmt.Errorf("dedup: field %q not found in entry", d.field)
		}
		return fmt.Sprintf("%v", val), nil
	}

	// Full-entry hash
	b, err := json.Marshal(entry)
	if err != nil {
		return "", fmt.Errorf("dedup: failed to marshal entry: %w", err)
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
