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
	ol.Insert("\x05")
	ol.Remove("\x05")
	if out := fmt.Sprint(ol.GetRange("", "\xff")); out != "[1 a aa b c]" {
		t.Errorf("Expected `[1 a aa b c]`, got `%v`", out)
	}
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
	if out := fmt.Sprint(ol.GetRange("1", "b")); out != "[1 a aa]" {
		t.Errorf("Expected `[1 a aa]`, got `%v`", out)
	}
}
