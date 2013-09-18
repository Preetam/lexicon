package lexicon

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/PreetamJinka/orderedlist"
)

var ErrConflict = errors.New("lexicon: version conflict")

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

func (lex *Lexicon) setHelper(key, value, version string) error {
	_, present := lex.hashmap[key]
	if !present {
		lex.hashmap[key] = &LexValue{
			Value:   value,
			version: generateVersion(""),
		}
		lex.list.Insert(key)
	} else {
		if version == "" {
			lex.hashmap[key].Value = value
			lex.hashmap[key].version = generateVersion("")
			return nil
		} else {
			if version != lex.hashmap[key].version {
				return ErrConflict
			} else {
				lex.hashmap[key].Value = value
				lex.hashmap[key].version = generateVersion(lex.hashmap[key].version)
			}
		}
	}

	return nil
}

// Set sets a key to a value.
func (lex *Lexicon) Set(params ...string) error {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	if len(params) == 2 {
		return lex.setHelper(params[0], params[1], "")
	} else {
		return lex.setHelper(params[0], params[1], params[2])
	}
}

func (lex *Lexicon) SetMany(kv map[string]string) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()
	for key := range kv {
		lex.setHelper(key, kv[key], "")
	}
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key string) (value string, version string) {
	return lex.hashmap[key].Value, lex.hashmap[key].version
}

// Remove deletes a key-value pair from the lexicon.
func (lex *Lexicon) Remove(key string) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()
	delete(lex.hashmap, key)
	lex.list.Remove(key)
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

func generateVersion(prev string) string {
	h := md5.New()
	io.WriteString(h, prev)
	return fmt.Sprintf("%x", h.Sum(nil))
}
