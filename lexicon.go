package lexicon

import (
	"math/rand"
	"sync"

	"github.com/PreetamJinka/gtreap"
)

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

// Lexicon is an ordered key-value store.
type Lexicon struct {
	treaps     map[uint64]*gtreap.Treap
	compare    func(a, b interface{}) int
	Hasher     func(interface{}) int
	version    uint64
	minVersion uint64

	mutex *sync.Mutex
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
	lex := &Lexicon{compare: compare, mutex: &sync.Mutex{}}
	lex.treaps = make(map[uint64]*gtreap.Treap, 10)
	lex.treaps[lex.version] = gtreap.NewTreap(lex.compareKeyValue)
	lex.Hasher = func(i interface{}) int { return rand.Int() }
	return lex
}

// Set sets a key to a value.
func (lex *Lexicon) Set(key, value interface{}) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	kv := KeyValue{Key: key, Value: value}
	lex.treaps[lex.version+1] = lex.treaps[lex.version].Delete(kv).Upsert(kv, lex.Hasher(key))
	lex.version++
}

// SetMany sets multiple key-value pairs.
func (lex *Lexicon) SetMany(kv map[interface{}]interface{}) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	for key, val := range kv {
		kv := KeyValue{Key: key, Value: val}
		if _, ok := lex.treaps[lex.version+1]; !ok {
			lex.treaps[lex.version+1] = lex.treaps[lex.version].Delete(kv).Upsert(kv, lex.Hasher(key))
		} else {
			lex.treaps[lex.version+1] = lex.treaps[lex.version+1].Delete(kv).Upsert(kv, lex.Hasher(key))
		}
	}

	lex.version++
}

// Get returns a value at the given key.
func (lex *Lexicon) Get(key interface{}, version ...uint64) interface{} {
	ver := lex.version

	if len(version) > 0 {
		ver = version[0]
	}

	val := lex.treaps[ver].Get(KeyValue{Key: key, Value: 0})
	if val != nil {
		return val.(KeyValue).Value
	}

	return nil
}

// Remove deletes a key-value pair.
func (lex *Lexicon) Remove(key interface{}) {
	lex.treaps[lex.version+1] = lex.treaps[lex.version].Delete(KeyValue{Key: key, Value: 0})
	lex.version++
}

// ClearRange removes a range of key-value pairs.
// The range is from [start, end).
func (lex *Lexicon) ClearRange(start, end interface{}) {
	lex.mutex.Lock()
	defer lex.mutex.Unlock()

	lex.treaps[lex.version].VisitAscend(KeyValue{Key: start, Value: 0}, func(i gtreap.Item) bool {
		if lex.compare(i.(KeyValue).Key, end) >= 0 {
			return false
		}
		if _, ok := lex.treaps[lex.version+1]; !ok {
			lex.treaps[lex.version+1] = lex.treaps[lex.version].Delete(i)
		} else {
			lex.treaps[lex.version+1] = lex.treaps[lex.version+1].Delete(i)
		}

		return true
	})

	lex.version++
}

// GetRange returns a slice of KeyValue structs.
// The range is from [start, end).
func (lex *Lexicon) GetRange(start, end interface{}, version ...uint64) (kv []KeyValue) {
	ver := lex.version

	if len(version) > 0 {
		ver = version[0]
	}

	kv = make([]KeyValue, 0, 10)

	lex.treaps[ver].VisitAscend(KeyValue{Key: start, Value: 0}, func(i gtreap.Item) bool {
		if lex.compare(i.(KeyValue).Key, end) >= 0 {
			return false
		}
		kv = append(kv, i.(KeyValue))
		return true
	})

	return
}
