package lexicon

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")

	if val, _, _ := lex.Get("foo"); val != "bar" {
		t.Errorf(`Expected "bar", got "%v".`, val)
	}

	lex.Set("foo", "baz")

	if val, _, _ := lex.Get("foo"); val != "baz" {
		t.Errorf(`Expected "baz", got "%v".`, val)
	}
}

func TestSetCollision(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")

	val, version, _ := lex.Get("foo")
	if val != "bar" {
		t.Errorf(`Expected "bar", got "%v".`, val)
	}

	lex.Set("foo", "baz", version)

	if val, _, _ := lex.Get("foo"); val != "baz" {
		t.Errorf(`Expected "baz", got "%v".`, val)
	}

	err := lex.Set("foo", "bazz", version)

	if err != ErrConflict {
		t.Errorf("Expected ErrConflict, got %v", err)
	}

	if val, _, _ := lex.Get("foo"); val != "baz" {
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

func TestSetMany(t *testing.T) {
	lex := New()
	pairs := map[string]string{
		"foo":    "bar",
		"foobar": "baz",
		"bar":    "foo",
		"a":      "1",
	}

	lex.SetMany(pairs)
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

	if val, _, err := lex.Get("foo"); err != ErrKeyNotPresent {
		t.Errorf(`Expected ErrKeyNotPresent, got value "%v".`, val)
	}
}

func BenchmarkBasicSetRemove(b *testing.B) {
	lex := New()

	for i := 0; i < b.N; i++ {
		lex.Set(string(i), "val")
	}

	for i := 0; i < b.N; i++ {
		lex.Remove(string(i))
	}
}
