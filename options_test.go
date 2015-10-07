package main

import (
	"reflect"
	"testing"
)

func TestMapVar(t *testing.T) {
	// empty
	s := &MapVar{}
	if len(s.Context) != 0 {
		t.Error("not empty")
	}

	// set one value
	s = &MapVar{}
	if err := s.Set("a=aye"); err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{"a": "aye"}
	if !reflect.DeepEqual(s.Context, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", s.Context)
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
	if !reflect.DeepEqual(s.Context, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", s.Context)
		t.Errorf("  want: %+v", want)
	}
	haveStr := s.String()
	if haveStr != "a=aye,b=bee" && haveStr != "b=bee,a=aye" {
		t.Error("String() invalid")
	}
}

func TestJsonVarOneSet(t *testing.T) {
	// empty
	j := &JsonVar{}
	if len(j.Context) != 0 {
		t.Error("not empty")
	}

	// set valid JsonVar
	value := `{"a":"aye","b":"bee"}`
	j = &JsonVar{}
	if err := j.Set(value); err != nil {
		t.Error(err)
	}
	want := map[string]interface{}{"a": "aye", "b": "bee"}
	if !reflect.DeepEqual(j.Context, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", j.Context)
		t.Errorf("  want: %+v", want)
	}

	// set invalid JsonVar
	value = `{a":"aye","b":"bee"}`
	j = &JsonVar{}
	if err := j.Set(value); err == nil {
		t.Error("no error")
	}
}

func TestJsonVarMultiSet(t *testing.T) {
	j := &JsonVar{}
	v1 := `{"a":"aye"}`
	v2 := `{"b":"bee"}`

	want := map[string]interface{}{"a": "aye"}
	if err := j.Set(v1); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(j.Context, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", j.Context)
		t.Errorf("  want: %+v", want)
	}

	want = map[string]interface{}{"a": "aye", "b": "bee"}
	if err := j.Set(v2); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(j.Context, want) {
		t.Error("value invalid")
		t.Errorf("  have: %+v", j.Context)
		t.Errorf("  want: %+v", want)
	}
}
