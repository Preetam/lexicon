package lexicon

import (
	"errors"
	"sync"

	"github.com/PreetamJinka/orderedlist"
)

var ErrConflict = errors.New("lexicon: version conflict")
var ErrKeyNotPresent = errors.New("lexicon: key not present")

type KeyValue struct {
	Key   string
	Value LexValue
}

type LexValue struct {
	Value   string
	version int
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

func (lex *Lexicon) setHelper(key string, value string, version int) error {
	_, present := lex.hashmap[key]
	if !present {
		lex.hashmap[key] = &LexValue{
			Value:   value,
			version: 0,
		}
		lex.list.Insert(key)
	} else {
		if version != lex.hashmap[key].version && version != -1 {
			return ErrConflict
		} else {
			lex.hashmap[key].Value = value
			lex.hashmap[key].version++
		}
	}

	return nil
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key string, value string, version int) error {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	return lex.setHelper(key, value, version)
}

func (lex *Lexicon) SetMany(kv map[string]string) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()
	for key := range kv {
		lex.setHelper(key, kv[key], -1)
	}
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key string) (string, int, error) {
	val, present := lex.hashmap[key]

	if !present {
		return "", -1, ErrKeyNotPresent
	}

	return val.Value, val.version, nil
}

// Remove deletes a key-value pair from the lexicon.
func (lex *Lexicon) Remove(key string) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()
	delete(lex.hashmap, key)
	lex.list.Remove(key)
}

func (lex *Lexicon) ClearRange(start string, end string) {
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
