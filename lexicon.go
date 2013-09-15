package lexicon

import (
	"sync"

	"github.com/PreetamJinka/orderedlist"
)

type KeyValue struct {
	Key   string
	Value LexValue
}

type LexValue struct {
	Value   string
	version string
}

// Lexicon is an ordered key-value store.
type Lexicon struct {
	list    *orderedlist.OrderedList
	hashmap map[string]*LexValue
	mutex   *sync.Mutex
}

// New returns an initialized lexicon.
func New() *Lexicon {
	return &Lexicon{
		list:    orderedlist.New(),
		hashmap: make(map[string]*LexValue),
		mutex:   &sync.Mutex{},
	}
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key, value string) {
	lex.mutex.Lock()
	_, present := lex.hashmap[key]
	if !present {
		lex.hashmap[key] = &LexValue{Value: value}
	}
	lex.list.Insert(key)
	lex.mutex.Unlock()
}

func (lex *Lexicon) SetMany(kv map[string]string) {
	lex.mutex.Lock()
	for key := range kv {
		_, present := lex.hashmap[key]
		if !present {
			lex.hashmap[key] = &LexValue{Value: kv[key]}
		} else {
			lex.hashmap[key].Value = kv[key]
		}
		lex.list.Insert(key)
	}
	lex.mutex.Unlock()
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key string) (value string, version string) {
	return lex.hashmap[key].Value, lex.hashmap[key].version
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
