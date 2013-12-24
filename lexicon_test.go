package lexicon

import (
	"testing"
)

func CompareStrings(a, b interface{}) (result int) {
	defer func() {
		if r := recover(); r != nil {
			// Log it?
		}
	}()

	aStr := a.(string)
	bStr := b.(string)

	if aStr > bStr {
		result = 1
	}

	if aStr < bStr {
		result = -1
	}

	return
}

func TestSetGet(t *testing.T) {
	lex := New(CompareStrings)
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
	lex := New(CompareStrings)
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
	lex := New(CompareStrings)
	pairs := map[interface{}]interface{}{
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
	lex := New(CompareStrings)
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
	lex := New(CompareStrings)

	if val := lex.Get("foo"); val != nil {
		t.Errorf(`Expected ErrKeyNotPresent, got value "%v".`, val)
	}
}

func BenchmarkBasicSetRemove(b *testing.B) {
	lex := New(CompareStrings)

	for i := 0; i < b.N; i++ {
		lex.Set(i, "val")
	}

	for i := 0; i < b.N; i++ {
		lex.Remove(i)
	}
}
