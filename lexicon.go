package lexicon

import (
	"./orderedlist"
	"sync"
)

type LexValue struct {
	value string
}

type LexKeyValue struct {
	Key   string
	Value string
}

type Lexicon struct {
	list    *orderedlist.OrderedList
	hashmap map[string]*LexValue
	mutex   *sync.Mutex
}

func New() *Lexicon {
	return &Lexicon{
		list:    orderedlist.New(),
		hashmap: make(map[string]*LexValue),
		mutex:   &sync.Mutex{},
	}
}

func (lex *Lexicon) Set(key string, value string) {
	lex.mutex.Lock()
	lex.hashmap[key] = &LexValue{value: value}
	lex.list.Insert(key)
	lex.mutex.Unlock()
}

func (lex *Lexicon) Get(key string) (value string) {
	lexvalue := lex.hashmap[key]
	value = lexvalue.value
	return
}

func (lex *Lexicon) Remove(key string) {
	delete(lex.hashmap, key)
	lex.list.Remove(key)
}

func (lex *Lexicon) GetRange(start string, end string) (kv []LexKeyValue) {
	kv = make([]LexKeyValue, 0)
	keys := lex.list.GetRange(start, end)

	for _, key := range keys {
		kv = append(kv, LexKeyValue{
			Key:   key,
			Value: lex.hashmap[key].value,
		})
	}
	return
}
