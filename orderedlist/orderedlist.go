package orderedlist

import (
	"fmt"
)

type OrderedList struct {
	slice []string
}

func New() *OrderedList {
	return &OrderedList{
		slice: make([]string, 0),
	}
}

func (l *OrderedList) Insert(key string) {
	// Empty list or greatest key
	if len(l.slice) == 0 {
		l.slice = append(l.slice, key)
		return
	}

	// Insert in O(n) time
	for index, val := range l.slice {
		if key < val {
			head := append(make([]string, 0), l.slice[:index]...)
			tail := append([]string{key}, l.slice[index:]...)
			l.slice = append(head, tail...)
			return
		}
	}

	l.slice = append(l.slice, key)
}

func (l *OrderedList) Remove(key string) {
	for index, val := range l.slice {
		if val == key {
			l.slice = append(l.slice[:index], l.slice[index+1:]...)
		}
	}
}

func (l *OrderedList) firstGreaterThanOrEqual(key string) (string, int) {
	index := 0
	for index, val := range l.slice {
		if val >= key {
			return val, index
		}
	}

	return l.slice[index], index
}

func (l *OrderedList) GetRange(start string, end string) (keys []string) {
	keys = make([]string, 0)
	_, index := l.firstGreaterThanOrEqual(start)
	for _, val := range l.slice[index:] {
		if val < end {
			keys = append(keys, val)
		}
	}
	return
}

func (l *OrderedList) Print() {
	fmt.Println(l.slice)
}
