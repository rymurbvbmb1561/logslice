package ratelimit

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_ZeroRate_AllowsAll(t *testing.T) {
	l := New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow() to return true for unlimited limiter at call %d", i)
		}
	}
}

func TestNew_NegativeRate_TreatedAsUnlimited(t *testing.T) {
	l := New(-5)
	if l.Rate() != 0 {
		t.Fatalf("expected rate 0 for negative input, got %d", l.Rate())
	}
	if !l.Allow() {
		t.Fatal("expected Allow() true for unlimited limiter")
	}
}

func TestAllow_RespectsRateWithinInterval(t *testing.T) {
	fixed := time.Now()
	l := New(3)
	l.clock = func() time.Time { return fixed }
	l.lastFill = fixed

	allowed := 0
	for i := 0; i < 6; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed within interval, got %d", allowed)
	}
}

func TestAllow_RefillsAfterInterval(t *testing.T) {
	base := time.Now()
	current := base
	l := New(2)
	l.clock = func() time.Time { return current }
	l.lastFill = base

	// drain the bucket
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected third call to be denied")
	}

	// advance one full second
	current = base.Add(time.Second)
	if !l.Allow() {
		t.Fatal("expected Allow() true after refill")
	}
}

func TestAllow_MultiplePeriods_CappedAtRate(t *testing.T) {
	base := time.Now()
	l := New(5)
	l.clock = func() time.Time { return base.Add(10 * time.Second) }
	l.lastFill = base
	l.tokens = 0

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	// tokens should have been capped at rate (5), not 10*5
	if allowed != 5 {
		t.Fatalf("expected 5 allowed after multi-period refill cap, got %d", allowed)
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	l := New(100)
	var count int64
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 20; j++ {
				if l.Allow() {
					atomic.AddInt64(&count, 1)
				}
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if count > 100 {
		t.Fatalf("expected at most 100 allowed, got %d", count)
	}
}
