// Package ratelimit provides a token-bucket style rate limiter for log
// pipeline output, allowing callers to cap the number of matched entries
// emitted per second.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls how many entries are allowed through per second.
type Limiter struct {
	mu       sync.Mutex
	rate     int           // max entries per interval
	interval time.Duration // refill interval
	tokens   int
	lastFill time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows at most rate entries per second.
// A rate of 0 disables limiting (all entries pass).
func New(rate int) *Limiter {
	if rate < 0 {
		rate = 0
	}
	return &Limiter{
		rate:     rate,
		interval: time.Second,
		tokens:   rate,
		lastFill: time.Now(),
		clock:    time.Now,
	}
}

// Allow reports whether the next entry should be allowed through.
// It is safe for concurrent use.
func (l *Limiter) Allow() bool {
	if l.rate == 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastFill)
	if elapsed >= l.interval {
		periods := int(elapsed / l.interval)
		l.tokens += periods * l.rate
		if l.tokens > l.rate {
			l.tokens = l.rate
		}
		l.lastFill = l.lastFill.Add(time.Duration(periods) * l.interval)
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// Rate returns the configured rate limit. Zero means unlimited.
func (l *Limiter) Rate() int {
	return l.rate
}
