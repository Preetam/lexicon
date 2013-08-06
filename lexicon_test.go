package lexicon

import (
	"fmt"
	"testing"
)

func TestSetGet(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")

	if val := lex.Get("foo"); val != "bar" {
		t.Errorf(`Expected "foo", got %v.`, val)
	}
}

func TestGetRange(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")
	lex.Set("foobar", "baz")
	lex.Set("bar", "foo")
	lex.Set("a", "1")

	kv := lex.GetRange("", "\xff")
	if res := fmt.Sprint(kv); res != "[{a 1} {bar foo} {foo bar} {foobar baz}]" {
		t.Errorf("Expected kv to be %v, got %v",
			"[{a 1} {bar foo} {foo bar} {foobar baz}]",
			res)
	}
}

func TestWriteRead(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")
	lex.Set("foobar", "baz")
	lex.Set("bar", "foo")
	lex.Set("a", "1")

	if err := lex.WriteToFile("/tmp/test.lex"); err != nil {
		t.Error(err)
	}

	lex, err := ReadFromFile("/tmp/test.lex")

	if err != nil {
		t.Error(err)
	}

	kv := lex.GetRange("", "\xff")
	if res := fmt.Sprint(kv); res != "[{a 1} {bar foo} {foo bar} {foobar baz}]" {
		t.Errorf("Expected kv to be %v, got %v",
			"[{a 1} {bar foo} {foo bar} {foobar baz}]",
			res)
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
