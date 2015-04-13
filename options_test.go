package main

import (
	"reflect"
	"testing"
)

func TestMapVar(t *testing.T) {
	// empty
	s := &MapVar{}
	if len(*s) != 0 {
		t.Error("not empty")
	}

	// set one value
	s = &MapVar{}
	if err := s.Set("a=aye"); err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{"a": "aye"}
	have := map[string]interface{}(*s)
	if !reflect.DeepEqual(have, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", have)
		t.Errorf("  want: %+v", want)
	}
	if s.String() != "a=aye" {
		t.Error("String() invalid")
	}

	// set multiple values
	s = &MapVar{}
	if err := s.Set("a=aye"); err != nil {
		t.Error(err)
	}
	if err := s.Set("b=bee"); err != nil {
		t.Error(err)
	}
	want = map[string]interface{}{"a": "aye", "b": "bee"}
	have = map[string]interface{}(*s)
	if !reflect.DeepEqual(have, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", have)
		t.Errorf("  want: %+v", want)
	}
	haveStr := s.String()
	if haveStr != "a=aye,b=bee" && haveStr != "b=bee,a=aye" {
		t.Error("String() invalid")
	}
}

func TestJsonVar(t *testing.T) {
	// empty
	j := &JsonVar{}
	if len(*j) != 0 {
		t.Error("not empty")
	}

	// set valid JsonVar
	value := `{"a":"aye","b":"bee"}`
	j = &JsonVar{}
	if err := j.Set(value); err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{"a": "aye", "b": "bee"}
	have := map[string]interface{}(*j)
	if !reflect.DeepEqual(have, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", have)
		t.Errorf("  want: %+v", want)
	}

	// set invalid JsonVar
	value = `{a":"aye","b":"bee"}`
	j = &JsonVar{}
	if err := j.Set(value); err == nil {
		t.Error("no error")
	}
}
