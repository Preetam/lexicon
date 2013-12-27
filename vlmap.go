package lexicon

/*
#cgo LDFLAGS: -lvlmap
#include <vlmap.h>

uint8_t*
helper(char* s) {
	return (uint8_t*)s;
}

char*
helper2(uint8_t* s) {
	return (char*)s;
}
*/
import "C"

import (
	"unsafe"
)

type vlMap struct {
	vlm *C.vlmap
}

func NewVlmap() *vlMap {
	return &vlMap{
		vlm: C.vlmap_create(),
	}
}

func (m *vlMap) Destroy() {
	C.vlmap_destroy(m.vlm)
}

func (m *vlMap) Version() uint64 {
	return uint64(C.vlmap_version(m.vlm))
}

func (m *vlMap) VersionInc() {
	C.vlmap_version_increment(m.vlm)
}

func (m *vlMap) Set(key string, value string) {
	C.vlmap_insert(m.vlm, C.vlmap_version(m.vlm), C.helper(C.CString(key)), C.int(len(key)),
		C.helper(C.CString(value)), C.int(len(value)))
}

func (m *vlMap) Get(key string, version uint64) (string, bool) {
	var val *C.uint8_t
	var vallen C.int
	getErr := C.vlmap_get(m.vlm, C.uint64_t(version), C.helper(C.CString(key)), C.int(len(key)), &val, &vallen)
	if getErr == 0 {
		value := C.GoStringN(C.helper2(val), vallen)
		C.free(unsafe.Pointer(val))
		return value, true
	} else {
		return "", false
	}
}

func (m *vlMap) Remove(key string) {
	C.vlmap_remove(m.vlm, C.vlmap_version(m.vlm), C.helper(C.CString(key)), C.int(len(key)))
}

func (m *vlMap) NewVlMapIterator(version uint64, startkey string, endkey string) *vlMapIterator {
	return &vlMapIterator{
		iter: C.vlmap_iterator_create(m.vlm, C.uint64_t(version),
			C.helper(C.CString(startkey)), C.int(len(startkey)),
			C.helper(C.CString(endkey)), C.int(len(endkey))),
	}
}

type vlMapIterator struct {
	iter *C.vlmap_iterator
}

func (i *vlMapIterator) Next() {
	next := C.vlmap_iterator_next(i.iter)
	if uintptr(unsafe.Pointer(next)) == 0 {
		C.vlmap_iterator_destroy(i.iter)
		i.iter = nil
		return
	}

	i.iter = next
}

func (i *vlMapIterator) Key() (string, bool) {
	var key *C.uint8_t
	var keylen C.int
	getErr := C.vlmap_iterator_get_key(i.iter, &key, &keylen)
	if getErr == 0 {
		keyStr := C.GoStringN(C.helper2(key), keylen)
		C.free(unsafe.Pointer(key))
		return keyStr, true
	} else {
		return "", false
	}
}

func (i *vlMapIterator) Value() (string, bool) {
	var val *C.uint8_t
	var vallen C.int
	getErr := C.vlmap_iterator_get_value(i.iter, &val, &vallen)
	if getErr == 0 {
		value := C.GoStringN(C.helper2(val), vallen)
		C.free(unsafe.Pointer(val))
		return value, true
	} else {
		return "", false
	}
}

func (m *vlMap) GetRange(version uint64, start string, end string) []KeyValue {
	ret := make([]KeyValue, 0, 10)

	i := m.NewVlMapIterator(version, start, end)
	for i.iter != nil {
		if key, ok := i.Key(); ok {
			if value, okVal := i.Value(); okVal {
				ret = append(ret, KeyValue{Key: key, Value: value})
			}
		}
		i.Next()
	}
	return ret
}

func (m *vlMap) ClearRange(version uint64, start string, end string) {
	i := m.NewVlMapIterator(version, start, end)
	for i.iter != nil {
		if _, ok := i.Key(); ok {
			next := C.vlmap_iterator_remove(i.iter)
			if uintptr(unsafe.Pointer(next)) == 0 {
				C.vlmap_iterator_destroy(i.iter)
				i.iter = nil
			}
		}
	}
}

func (i *vlMapIterator) Destroy() {
	C.vlmap_iterator_destroy(i.iter)
	i.iter = nil
}
