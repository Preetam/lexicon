package lexicon

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	lex := New()
	lex.Set("foo", "bar")

	if val := lex.Get("foo"); val != "bar" {
		t.Errorf(`Expected "foo", got %v.`, val)
	}
}
