// Package sampler provides log line sampling functionality for logslice.
// It supports rate-based sampling (every Nth line) and random sampling
// to reduce output volume when working with large log files.
package sampler

import (
	"math/rand"
)

// Sampler decides whether a given log line should be included in output.
type Sampler struct {
	rate   int
	random float64
	count  int
	rng    *rand.Rand
}

// Option configures a Sampler.
type Option func(*Sampler)

// WithRate configures the sampler to keep every Nth line (1 = keep all).
func WithRate(n int) Option {
	return func(s *Sampler) {
		if n > 0 {
			s.rate = n
		}
	}
}

// WithRandom configures the sampler to keep each line with probability p (0.0–1.0).
func WithRandom(p float64) Option {
	return func(s *Sampler) {
		if p >= 0.0 && p <= 1.0 {
			s.random = p
		}
	}
}

// New creates a new Sampler with the given options.
// By default all lines are kept (rate=1, random=1.0).
func New(seed int64, opts ...Option) *Sampler {
	s := &Sampler{
		rate:   1,
		random: 1.0,
		rng:    rand.New(rand.NewSource(seed)),
	}
	for _, o := opts {
		o(s)
	}
	return s
}

// Keep returns true if the current line should be included in output.
// Both rate and random filters must pass for a line to be kept.
func (s *Sampler) Keep() bool {
	s.count++

	// Rate filter: keep only every Nth line.
	if s.rate > 1 && s.count%s.rate != 0 {
		return false
	}

	// Random filter: keep with probability s.random.
	if s.random < 1.0 && s.rng.Float64() >= s.random {
		return false
	}

	return true
}

// Reset resets the internal line counter.
func (s *Sampler) Reset() {
	s.count = 0
}
