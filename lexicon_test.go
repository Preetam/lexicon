package lexicon

import (
	"testing"

	"github.com/PreetamJinka/orderedlist"
)

type ComparableString string

func (cs ComparableString) Compare(c orderedlist.Comparable) int {
	if cs > c.(ComparableString) {
		return 1
	}
	if cs < c.(ComparableString) {
		return -1
	}
	return 0
}

func TestSetGet(t *testing.T) {
	lex := New()
	lex.Set(ComparableString("foo"), "bar")

	if val, _ := lex.Get(ComparableString("foo")); val != "bar" {
		t.Errorf(`Expected "bar", got "%v".`, val)
	}

	lex.Set(ComparableString("foo"), "baz")

	if val, _ := lex.Get(ComparableString("foo")); val != "baz" {
		t.Errorf(`Expected "baz", got "%v".`, val)
	}
}

func TestGetRange(t *testing.T) {
	lex := New()
	lex.Set(ComparableString("foo"), "bar")
	lex.Set(ComparableString("foobar"), "baz")
	lex.Set(ComparableString("bar"), "foo")
	lex.Set(ComparableString("a"), "1")

	kv := lex.GetRange(ComparableString(""), ComparableString("\xff"))
	if len(kv) != 4 {
		t.Errorf("Expected 4 results, got %d", len(kv))
	}
}

func TestSetMany(t *testing.T) {
	lex := New()
	pairs := map[orderedlist.Comparable]interface{}{
		ComparableString("foo"):    "bar",
		ComparableString("foobar"): "baz",
		ComparableString("bar"):    "foo",
		ComparableString("a"):      "1",
	}

	lex.SetMany(pairs)
	kv := lex.GetRange(ComparableString(""), ComparableString("\xff"))
	if len(kv) != 4 {
		t.Errorf("Expected 4 results, got %d", len(kv))
	}
}

func TestClearRange(t *testing.T) {
	lex := New()
	lex.Set(ComparableString("foo"), "bar")
	lex.Set(ComparableString("foobar"), "baz")
	lex.Set(ComparableString("bar"), "foo")
	lex.Set(ComparableString("a"), "1")

	lex.ClearRange(ComparableString("foo"), ComparableString("foobar\xff"))

	kv := lex.GetRange(ComparableString(""), ComparableString("\xff"))
	if len(kv) != 2 {
		t.Errorf("Expected 2 results, got %d", len(kv))
	}
}

func TestMissingKey(t *testing.T) {
	lex := New()

	if val, err := lex.Get(ComparableString("foo")); err != ErrKeyNotPresent {
		t.Errorf(`Expected ErrKeyNotPresent, got value "%v".`, val)
	}
}

func BenchmarkBasicSetRemove(b *testing.B) {
	lex := New()

	for i := 0; i < b.N; i++ {
		lex.Set(ComparableString(i), "val")
	}

	for i := 0; i < b.N; i++ {
		lex.Remove(ComparableString(i))
	}
}
