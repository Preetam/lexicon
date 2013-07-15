package orderedlist

import (
	"fmt"
	"testing"
)

func TestSetGet(t *testing.T) {
	ol := New()

	ol.Insert("c")
	ol.Insert("a")
	ol.Insert("b")
	ol.Insert("aa")
	ol.Insert("1")
	ol.Insert("\xfa")
	ol.Remove("\xfa")
	ol.Print()
}

func TestGetRange(t *testing.T) {
	ol := New()

	ol.Insert("c")
	ol.Insert("a")
	ol.Insert("b")
	ol.Insert("aa")
	ol.Insert("1")
	ol.Insert("\x05")
	ol.Remove("\x05")
	fmt.Println(ol.GetRange("1", "b"))
	fmt.Println("-------")
	fmt.Println(ol.GetRange("", "\xff"))
}
