package main

import (
	"reflect"
	"testing"
)

func TestContext(t *testing.T) {
	// empty context
	c := &Context{}
	have := c.Map()
	if len(have) != 0 {
		t.Error("map is not empty")
	}

	// single merge
	in := map[string]interface{}{"a": "aye", "b": "bee"}
	want := map[string]interface{}{"a": "aye", "b": "bee"}
	c = &Context{}
	c.Update(in)
	have = c.Map()
	if !reflect.DeepEqual(have, want) {
		t.Error("maps are not equal")
	}

	// overwrite, depth of 1
	in = map[string]interface{}{"a": "aye", "b": "bee"}
	update := map[string]interface{}{"a": "eh"}
	want = map[string]interface{}{"a": "eh", "b": "bee"}
	c = &Context{}
	c.Update(in)
	c.Update(update)
	have = c.Map()
	if !reflect.DeepEqual(have, want) {
		t.Error("maps are not equal")
	}

	// overwrite, let's go deeper
	in = map[string]interface{}{
		"a": []interface{}{"aye", "eh"},
		"b": "bee",
		"c": "see",
	}
	update = map[string]interface{}{
		"c": []interface{}{"see", "sea"},
	}
	want = map[string]interface{}{
		"a": []interface{}{"aye", "eh"},
		"b": "bee",
		"c": []interface{}{"see", "sea"},
	}
	c = &Context{}
	c.Update(in)
	c.Update(update)
	have = c.Map()
	if !reflect.DeepEqual(have, want) {
		t.Error("maps are not equal")
	}
}
