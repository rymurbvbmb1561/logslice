package sampler

import (
	"testing"
)

func TestNew_DefaultKeepsAll(t *testing.T) {
	s := New(42)
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Fatalf("default sampler dropped line %d", i+1)
		}
	}
}

func TestWithRate_KeepsEveryNth(t *testing.T) {
	s := New(42, WithRate(3))
	kept := 0
	for i := 0; i < 9; i++ {
		if s.Keep() {
			kept++
		}
	}
	if kept != 3 {
		t.Errorf("expected 3 kept lines, got %d", kept)
	}
}

func TestWithRate_One_KeepsAll(t *testing.T) {
	s := New(42, WithRate(1))
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Fatalf("rate=1 sampler dropped line %d", i+1)
		}
	}
}

func TestWithRate_Zero_IgnoredDefaultsToOne(t *testing.T) {
	s := New(42, WithRate(0))
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Fatalf("rate=0 (ignored) sampler dropped line %d", i+1)
		}
	}
}

func TestWithRandom_Zero_DropsAll(t *testing.T) {
	s := New(42, WithRandom(0.0))
	for i := 0; i < 50; i++ {
		if s.Keep() {
			t.Fatalf("random=0.0 sampler kept line %d", i+1)
		}
	}
}

func TestWithRandom_One_KeepsAll(t *testing.T) {
	s := New(42, WithRandom(1.0))
	for i := 0; i < 50; i++ {
		if !s.Keep() {
			t.Fatalf("random=1.0 sampler dropped line %d", i+1)
		}
	}
}

func TestWithRandom_Half_ApproxHalf(t *testing.T) {
	s := New(99, WithRandom(0.5))
	kept := 0
	total := 1000
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}
	// Allow ±15% tolerance around 50%.
	if kept < 350 || kept > 650 {
		t.Errorf("expected ~500 kept lines, got %d", kept)
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s := New(42, WithRate(2))
	// Consume 4 lines: lines 2 and 4 kept.
	kept := 0
	for i := 0; i < 4; i++ {
		if s.Keep() {
			kept++
		}
	}
	if kept != 2 {
		t.Fatalf("pre-reset: expected 2 kept, got %d", kept)
	}
	s.Reset()
	// After reset line 2 should be kept again.
	s.Keep() // line 1 — dropped
	if !s.Keep() { // line 2 — kept
		t.Error("post-reset: expected line 2 to be kept")
	}
}
