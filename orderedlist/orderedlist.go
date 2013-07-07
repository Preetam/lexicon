package orderedlist

import (
	"container/list"
	"fmt"
)

type OrderedList struct {
	linkedlist *list.List
}

func New() *OrderedList {
	return &OrderedList{
		linkedlist: list.New(),
	}
}

func (l *OrderedList) Insert(key string) {
	// Empty list or greatest key
	if l.linkedlist.Len() == 0 || l.linkedlist.Back().Value.(string) < key {
		l.linkedlist.PushBack(key)
		return
	}

	// Insert in O(n) time
	for e := l.linkedlist.Front(); e != nil; e = e.Next() {
		if e.Value.(string) > key {
			l.linkedlist.InsertBefore(key, e)
			return
		}
	}
}

func (l *OrderedList) Remove(key string) {
	for e := l.linkedlist.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == key {
			l.linkedlist.Remove(e)
			return
		}
	}
}

func (l *OrderedList) firstGreaterThanOrEqual(key string) *list.Element {
	elem := l.linkedlist.Front()
	for e := elem; e != nil; e = e.Next() {
		if e.Value.(string) >= key {
			return e
		}
	}

	return elem
}

func (l *OrderedList) GetRange(start string, end string) (keys []string) {
	keys = make([]string, 0)
	startElem := l.firstGreaterThanOrEqual(start)
	for e := startElem; e != nil; e = e.Next() {
		if e.Value.(string) < end {
			keys = append(keys, e.Value.(string))
		}
	}
	return
}

func (l *OrderedList) Print() {
	for e := l.linkedlist.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
