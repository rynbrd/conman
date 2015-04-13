package main

import (
	"reflect"
	"testing"
)

func envsEqual(a, b []string) bool {
	makeSet := func(items []string) map[string]struct{} {
		set := make(map[string]struct{}, len(items))
		for _, item := range items {
			set[item] = struct{}{}
		}
		return set
	}
	return reflect.DeepEqual(makeSet(a), makeSet(b))
}

func TestEnviron(t *testing.T) {
	// empty environ
	e := &Environ{}
	v := e.Values()
	if len(v) != 0 {
		t.Error("environ is not empty")
	}

	// single add
	in := []string{"A=aye", "B=bee"}
	want := []string{"A=aye", "B=bee"}
	e = &Environ{}
	v = e.Values()
	e.Load(in)
	have := e.Values()
	if !envsEqual(have, want) {
		t.Error("environ values not equal")
		t.Errorf("  have: %+v\n", have)
		t.Errorf("  want: %+v\n", want)
	}

	// add with duplicates
	in = []string{"A=aye", "B=bee", "A=eh"}
	want = []string{"A=aye", "B=bee"}
	e = &Environ{}
	v = e.Values()
	e.Load(in)
	have = e.Values()
	if !envsEqual(have, want) {
		t.Error("environ values not equal")
		t.Errorf("  have: %+v\n", have)
		t.Errorf("  want: %+v\n", want)
	}

	// replace
	in = []string{"A=aye", "B=bee"}
	update := []string{"A=eh", "B=bee"}
	want = []string{"A=eh", "B=bee"}
	e = &Environ{}
	v = e.Values()
	e.Load(in)
	e.Load(update)
	have = e.Values()
	if !envsEqual(have, want) {
		t.Error("environ values not equal")
		t.Errorf("  have: %+v\n", have)
		t.Errorf("  want: %+v\n", want)
	}
}
