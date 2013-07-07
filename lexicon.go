package lexicon

import (
	"github.com/huandu/skiplist"
	"sync"
)

type LexValue struct {
	value string
}

type Lexicon struct {
	list    *skiplist.SkipList
	hashmap map[string]*LexValue
	mutex   *sync.Mutex
}

func New() *Lexicon {
	return &Lexicon{
		list:    skiplist.New(skiplist.StringAsc),
		hashmap: make(map[string]*LexValue),
		mutex:   &sync.Mutex{},
	}
}

func (lex *Lexicon) Set(key string, value string) {
	lex.mutex.Lock()
	lex.hashmap[key] = &LexValue{value: value}
	lex.list.Set(key, lex.hashmap[key])
	lex.mutex.Unlock()
}

func (lex *Lexicon) Get(key string) (value string) {
	lexvalue := lex.hashmap[key]
	value = lexvalue.value
	return
}
