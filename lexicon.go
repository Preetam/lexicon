package lexicon

import (
	"errors"
	"math/rand"

	"github.com/PreetamJinka/gtreap"
)

var ErrKeyNotPresent = errors.New("lexicon: key not present")

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// Lexicon is an ordered key-value store.
type Lexicon struct {
	treap   *gtreap.Treap
	compare func(a, b interface{}) int
}

func (lex *Lexicon) compareKeyValue(a, b interface{}) int {
	akv, ok1 := a.(KeyValue)
	bkv, ok2 := b.(KeyValue)
	if !(ok1 && ok2) {
		return 0
	}

	return lex.compare(akv.Key, bkv.Key)
}

// New returns an initialized lexicon.
func New(compare func(a, b interface{}) int) *Lexicon {
	lex := &Lexicon{compare: compare}
	lex.treap = gtreap.NewTreap(lex.compareKeyValue)
	return lex
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key, value interface{}) {
	kv := KeyValue{Key: key, Value: value}
	lex.treap = lex.treap.Delete(kv).Upsert(kv, rand.Int())
}

// SetMany sets multiple key-value pairs.
func (lex *Lexicon) SetMany(kv map[interface{}]interface{}) {
	for key, val := range kv {
		lex.Set(key, val)
	}
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key interface{}) interface{} {
	val := lex.treap.Get(KeyValue{Key: key, Value: 0})
	if val != nil {
		return val.(KeyValue).Value
	}

	return nil
}

// Remove deletes a key-value pair.
func (lex *Lexicon) Remove(key interface{}) {
	lex.treap = lex.treap.Delete(KeyValue{Key: key, Value: 0})
}

// ClearRange removes a range of key-value pairs.
// The range is from [start, end).
func (lex *Lexicon) ClearRange(start, end interface{}) {
	lex.treap.VisitAscend(KeyValue{Key: start, Value: 0}, func(i gtreap.Item) bool {
		if lex.compare(i.(KeyValue).Key, end) >= 0 {
			return false
		}
		lex.Remove(i.(KeyValue).Key)
		return true
	})
}

// GetRange returns a slice of KeyValue structs.
// The range is from [start, end).
func (lex *Lexicon) GetRange(start, end interface{}) (kv []KeyValue) {
	kv = make([]KeyValue, 0, 10)

	lex.treap.VisitAscend(KeyValue{Key: start, Value: 0}, func(i gtreap.Item) bool {
		if lex.compare(i.(KeyValue).Key, end) >= 0 {
			return false
		}
		kv = append(kv, i.(KeyValue))
		return true
	})

	return
}
