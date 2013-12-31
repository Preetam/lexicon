package lexicon

import (
	"sync"
)

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// Lexicon is an ordered key-value store.
type Lexicon struct {
	m     *vlMap
	mutex *sync.Mutex
}

// New returns an initialized lexicon.
func New() *Lexicon {
	return &Lexicon{
		m:     NewVlmap(),
		mutex: &sync.Mutex{},
	}
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key, value string) {
	lex.m.Set(key, value)
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key string) string {
	val, _ := lex.m.Get(key, 1)
	return val
}

// Remove deletes a key-value pair.
func (lex *Lexicon) Remove(key string) {
	lex.m.Remove(key)
}

// ClearRange removes a range of key-value pairs.
// The range is from [start, end).
func (lex *Lexicon) ClearRange(start, end string) {
	lex.m.ClearRange(1, start, end)
}

// GetRange returns a slice of KeyValue structs.
// The range is from [start, end).
func (lex *Lexicon) GetRange(start, end string) (kv []KeyValue) {
	return lex.m.GetRange(1, start, end)
}
