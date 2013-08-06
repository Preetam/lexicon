package lexicon

import (
	"encoding/gob"
	"github.com/PreetamJinka/orderedlist"
	"os"
	"sync"
)

type LexValue struct {
	value interface{}
}

type LexKeyValue struct {
	Key   string
	Value interface{}
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
func (lex *Lexicon) Set(key string, value interface{}) {
	lex.mutex.Lock()
	lex.hashmap[key] = &LexValue{value: value}
	lex.list.Insert(key)
	lex.mutex.Unlock()
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key string) (value interface{}) {
	lexvalue := lex.hashmap[key]
	value = lexvalue.value
	return
}

// Remove deletes a key-value pair from the lexicon.
func (lex *Lexicon) Remove(key string) {
	lex.mutex.Lock()
	delete(lex.hashmap, key)
	lex.list.Remove(key)
	lex.mutex.Unlock()
}

// GetRange returns a slice of LexKeyValue structs.
// The range is from [start, end).
func (lex *Lexicon) GetRange(start string, end string) (kv []LexKeyValue) {
	kv = make([]LexKeyValue, 0)

	lex.mutex.Lock()
	keys := lex.list.GetRange(start, end)
	lex.mutex.Unlock()

	for _, key := range keys {
		kv = append(kv, LexKeyValue{
			Key:   key,
			Value: lex.hashmap[key].value,
		})
	}
	return
}

// WriteToFile writes a LexKeyValue slice to fileName.
func (lex *Lexicon) WriteToFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(lex.GetRange("", "\xff"))
}

func readKVFromFile(fileName string) ([]LexKeyValue, error) {
	var lex []LexKeyValue
	file, err := os.Open(fileName)
	if err != nil {
		return lex, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&lex)
	return lex, err
}

// ReadFromFile returns a new Lexicon generated from
// key-value pairs in fileName.
func ReadFromFile(fileName string) (*Lexicon, error) {
	kvSlice, err := readKVFromFile(fileName)
	if err != nil {
		return nil, err
	}

	lex := New()

	for _, kv := range kvSlice {
		lex.Set(kv.Key, kv.Value)
	}

	return lex, nil
}
