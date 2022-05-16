package main

import (
	"container/list"
	"reflect"
	"testing"
)

func Equal(x any, y any) bool {

	return reflect.DeepEqual(x, y)

}

func TestStringify(t *testing.T) {

	got := Stringify(DummyInfo{
		N: 6,
		M: 9,
	})
	want := "{\"N\":6,\"M\":9}"

	if got != want {
		t.Errorf("Test stringify failed!")
	}

}

func TestDestringify(t *testing.T) {

	got1 := Destringify("{\"N\":6,\"M\":9}")
	want1 := DummyInfo{
		N: 6,
		M: 9,
	}

	if !Equal(got1, want1) {
		t.Errorf("Test single element destringify failed!")
		return
	}

	got2 := Destringify("[{\"N\":5,\"M\":10}, {\"N\":6,\"M\":9}]")
	want2 := list.New()
	want2.PushBack(DummyInfo{
		N: 5,
		M: 10,
	})
	want2.PushBack(DummyInfo{
		N: 6,
		M: 9,
	})

	got2_list, ok := got2.(list.List)
	if !ok {
		t.Errorf("Test list of elements destringify failed!")
		return
	}

	it1 := got2_list.Front()
	for it := want2.Front(); it != nil; it = it.Next() {
		if it1 == nil {
			t.Errorf("Test list of elements destringify failed!")
			return
		}
		if !Equal(it.Value, it1.Value) {
			t.Errorf("Test list of elements destringify failed!")
			return
		}
		it1 = it1.Next()
	}

}
