// Package orderedlist is a basic wrapper around container/list.
package orderedlist

import (
	"container/list"
	"fmt"
)

// OrderedList is an ordered linked list.
type OrderedList struct {
	linkedlist *list.List
}

// New returns an initialized OrderedList.
func New() *OrderedList {
	return &OrderedList{
		linkedlist: list.New(),
	}
}

// Insert inserts a key string into the ordered list.
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

// Remove removes a key from the ordered list.
func (l *OrderedList) Remove(key string) {
	for e := l.linkedlist.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == key {
			l.linkedlist.Remove(e)
			return
		}
	}
}

// firstGreaterThanOrEqual returns the first Element
// greater-than or equal-to the given key.
func (l *OrderedList) firstGreaterThanOrEqual(key string) *list.Element {
	elem := l.linkedlist.Front()
	for e := elem; e != nil; e = e.Next() {
		if e.Value.(string) >= key {
			return e
		}
	}

	return elem
}

// GetRange returns a slice of strings from [start, end).
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

// Print prints the values stored in the list.
func (l *OrderedList) Print() {
	for e := l.linkedlist.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
