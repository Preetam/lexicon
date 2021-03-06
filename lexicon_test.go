package lexicon

import (
	"fmt"
	"testing"
)

func TestSetGet(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")

	if val := lex.Get("foo"); val != "bar" {
		t.Errorf(`Expected "bar", got "%v".`, val)
	}

	lex.Set("foo", "baz")

	if val := lex.Get("foo"); val != "baz" {
		t.Errorf(`Expected "baz", got "%v".`, val)
	}
}

func TestHash(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")

	if val := lex.Get("foo"); val != "bar" {
		t.Errorf(`Expected "bar", got "%v".`, val)
	}

	lex.Set("foo", "baz")

	if val := lex.Get("foo"); val != "baz" {
		t.Errorf(`Expected "baz", got "%v".`, val)
	}
}

func TestGetRange(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")
	lex.Set("foobar", "baz")
	lex.Set("bar", "foo")
	lex.Set("a", "1")

	kv := lex.GetRange("", "\xff")
	if len(kv) != 4 {
		t.Errorf("Expected 4 results, got %d", len(kv))
	}
}

func TestClearRange(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")
	lex.Set("foobar", "baz")
	lex.Set("bar", "foo")
	lex.Set("a", "1")

	lex.ClearRange("foo", "foobar\xff")

	kv := lex.GetRange("", "\xff")
	if len(kv) != 2 {
		t.Errorf("Expected 2 results, got %d", len(kv))
	}
}

func TestMissingKey(t *testing.T) {
	lex := New()

	if val := lex.Get("foo"); val != "" {
		t.Errorf(`Expected ErrKeyNotPresent, got value "%v".`, val)
	}
}

func BenchmarkBasicSet(b *testing.B) {
	lex := New()

	for i := 0; i < b.N; i++ {
		lex.Set(fmt.Sprint(i), "val")
	}
}
