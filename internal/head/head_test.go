package head_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/head"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"` + msg + `"}`,
		Fields:    map[string]interface{}{"msg": msg},
	}
}

func TestApply_ZeroMax_AllowsAll(t *testing.T) {
	l := head.New(head.WithMax(0))
	for i := 0; i < 100; i++ {
		_, ok := l.Apply(makeEntry("x"))
		if !ok {
			t.Fatalf("expected entry %d to pass through", i)
		}
	}
}

func TestApply_NegativeMax_AllowsAll(t *testing.T) {
	l := head.New(head.WithMax(-5))
	for i := 0; i < 10; i++ {
		_, ok := l.Apply(makeEntry("x"))
		if !ok {
			t.Fatalf("expected entry %d to pass through", i)
		}
	}
}

func TestApply_LimitsToMax(t *testing.T) {
	l := head.New(head.WithMax(3))

	for i := 0; i < 3; i++ {
		_, ok := l.Apply(makeEntry("x"))
		if !ok {
			t.Fatalf("expected entry %d to pass through", i)
		}
	}

	_, ok := l.Apply(makeEntry("x"))
	if ok {
		t.Fatal("expected entry after limit to be rejected")
	}
}

func TestDone_FalseBeforeLimit(t *testing.T) {
	l := head.New(head.WithMax(5))
	l.Apply(makeEntry("x"))
	if l.Done() {
		t.Fatal("expected Done() to be false before limit reached")
	}
}

func TestDone_TrueAfterLimit(t *testing.T) {
	l := head.New(head.WithMax(1))
	l.Apply(makeEntry("x"))
	if !l.Done() {
		t.Fatal("expected Done() to be true after limit reached")
	}
}

func TestReset_AllowsReuse(t *testing.T) {
	l := head.New(head.WithMax(2))
	l.Apply(makeEntry("a"))
	l.Apply(makeEntry("b"))

	if !l.Done() {
		t.Fatal("expected Done() true before reset")
	}

	l.Reset()

	if l.Done() {
		t.Fatal("expected Done() false after reset")
	}

	_, ok := l.Apply(makeEntry("c"))
	if !ok {
		t.Fatal("expected entry to pass after reset")
	}
}
