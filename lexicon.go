package lexicon

import (
	"github.com/PreetamJinka/orderedlist"
	"sync"
)

type KeyValue struct {
	Key   string
	Value interface{}
}

// Lexicon is an ordered key-value store.
type Lexicon struct {
	list    *orderedlist.OrderedList
	hashmap map[string]*interface{}
	mutex   *sync.Mutex
}

// New returns an initialized lexicon.
func New() *Lexicon {
	return &Lexicon{
		list:    orderedlist.New(),
		hashmap: make(map[string]*interface{}),
		mutex:   &sync.Mutex{},
	}
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key string, value interface{}) {
	lex.mutex.Lock()
	lex.hashmap[key] = &value
	lex.list.Insert(key)
	lex.mutex.Unlock()
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key string) (value interface{}) {
	return *lex.hashmap[key]
}

// Remove deletes a key-value pair from the lexicon.
func (lex *Lexicon) Remove(key string) {
	lex.mutex.Lock()
	delete(lex.hashmap, key)
	lex.list.Remove(key)
	lex.mutex.Unlock()
}

// GetRange returns a slice of KeyValue structs.
// The range is from [start, end).
func (lex *Lexicon) GetRange(start string, end string) (kv []KeyValue) {
	kv = make([]KeyValue, 0)

	lex.mutex.Lock()
	keys := lex.list.GetRange(start, end)
	lex.mutex.Unlock()

	for _, key := range keys {
		kv = append(kv, KeyValue{
			Key:   key,
			Value: *lex.hashmap[key],
		})
	}
	return
}
