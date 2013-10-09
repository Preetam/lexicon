package lexicon

import (
	"errors"
	"sync"

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

var ErrKeyNotPresent = errors.New("lexicon: key not present")

type KeyValue struct {
	Key   orderedlist.Comparable
	Value interface{}
}

// Lexicon is an ordered key-value store.
type Lexicon struct {
	list    *orderedlist.OrderedList
	hashmap map[orderedlist.Comparable]interface{}
	mutex   *sync.Mutex
}

// New returns an initialized lexicon.
func New() *Lexicon {
	return &Lexicon{
		list:    orderedlist.New(),
		hashmap: make(map[orderedlist.Comparable]interface{}),
		mutex:   &sync.Mutex{},
	}
}

func (lex *Lexicon) setHelper(key orderedlist.Comparable, value interface{}) error {
	_, present := lex.hashmap[key]

	if !present {
		lex.hashmap[key] = value
		lex.list.Insert(key)
	} else {
		lex.hashmap[key] = value
	}

	return nil
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key orderedlist.Comparable, value interface{}) error {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	return lex.setHelper(key, value)
}

func (lex *Lexicon) SetMany(kv map[orderedlist.Comparable]interface{}) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()
	for key := range kv {
		lex.setHelper(key, kv[key])
	}
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key orderedlist.Comparable) (interface{}, error) {
	val, present := lex.hashmap[key]

	if !present {
		return "", ErrKeyNotPresent
	}

	return val, nil
}

// Remove deletes a key-value pair from the lexicon.
func (lex *Lexicon) Remove(key orderedlist.Comparable) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()
	delete(lex.hashmap, key)
	lex.list.Remove(key)
}

func (lex *Lexicon) ClearRange(start orderedlist.Comparable, end orderedlist.Comparable) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	keys := lex.list.GetRange(start, end)

	for _, key := range keys {
		delete(lex.hashmap, key)
		lex.list.Remove(key)
	}
}

// GetRange returns a slice of KeyValue structs.
// The range is from [start, end).
func (lex *Lexicon) GetRange(start orderedlist.Comparable, end orderedlist.Comparable) (kv []KeyValue) {
	kv = make([]KeyValue, 0)

	lex.mutex.Lock()
	keys := lex.list.GetRange(start, end)
	lex.mutex.Unlock()

	for _, key := range keys {
		kv = append(kv, KeyValue{
			Key:   key,
			Value: lex.hashmap[key],
		})
	}
	return
}
