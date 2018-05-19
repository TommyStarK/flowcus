package orderedmap

import (
	"errors"
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	orderedMap := NewOrderedMap()

	if len(orderedMap.keys) != 0 {
		t.Errorf("keys' length should equal to 0")
	}

	if len(orderedMap.mmap) != 0 {
		t.Errorf("map's length should be equal to 0")
	}

	if orderedMap.mutex == nil {
		t.Errorf("mutex should be of type *sync.Mutex ")
	}
}

func TestOrderedMapUseCases(t *testing.T) {
	orderedMap := NewOrderedMap()

	if should_be_nil := orderedMap.Get("should_be_nil"); should_be_nil != nil {
		t.Errorf("trying to get an element with a key which does not exist should return nil")
	}

	orderedMap.Set("toto", 1)
	orderedMap.Set(2, "titi")
	orderedMap.Set(3.4, "tata")
	orderedMap.Set(errors.New("new error"), true)
	orderedMap.Set("func", func(s string) int {
		t.Logf("print from func stored in ordred map: %s", s)
		return len(s)
	})

	if len := orderedMap.Len(); len != 5 {
		t.Errorf("5 elements have been added to our map, its length should be 5")
	}

	if elem := orderedMap.Get(3.4); elem == nil {
		t.Errorf("element should not be nil")
	} else {
		switch elem.(type) {
		case string:
			t.Log("elem type is valid")
		default:
			t.Errorf("elem should be of type string")
		}
	}

	orderedMap.Delete(3.4)
	orderedMap.Delete("invalid key")

	if len := orderedMap.Len(); len != 4 {
		t.Errorf("orderedMap's length should be equal to 4")
	}

	orderedMap.Set("toto", 2)
	item := orderedMap.Get("toto")
	if item == nil || item.(int) != 2 {
		t.Errorf("item should not be nil but of type int and equal to 2")
	}

	for _, key := range orderedMap.Keys() {
		elem := orderedMap.Get(key)
		switch elem.(type) {
		case func(string) int:
			if res := elem.(func(string) int)("debug"); res != 5 {
				t.Errorf("res should be equal to 5")
			}
		}
	}
}
