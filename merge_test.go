package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMerge(t *testing.T) {
	for _, tuple := range []struct {
		src         string
		dst         string
		expected    string
		arrayAppend bool
	}{
		{
			src:         `{}`,
			dst:         `{}`,
			expected:    `{}`,
			arrayAppend: false,
		},
		{
			src:         `{"b":2}`,
			dst:         `{"a":1}`,
			expected:    `{"a":1,"b":2}`,
			arrayAppend: false,
		},
		{
			src:         `{"a":0}`,
			dst:         `{"a":1}`,
			expected:    `{"a":0}`,
			arrayAppend: false,
		},
		{
			src:         `{"a":{       "y":2}}`,
			dst:         `{"a":{"x":1       }}`,
			expected:    `{"a":{"x":1, "y":2}}`,
			arrayAppend: false,
		},
		{
			src:         `{"a":{"x":2}}`,
			dst:         `{"a":{"x":1}}`,
			expected:    `{"a":{"x":2}}`,
			arrayAppend: false,
		},
		{
			src:         `{"a":{       "y":7, "z":8}}`,
			dst:         `{"a":{"x":1, "y":2       }}`,
			expected:    `{"a":{"x":1, "y":7, "z":8}}`,
			arrayAppend: false,
		},
		{
			src:         `{"1": { "b":1, "2": { "3": {         "b":3, "n":[1,2]} }        }}`,
			dst:         `{"1": {        "2": { "3": {"a":"A",        "n":"xxx"} }, "a":3 }}`,
			expected:    `{"1": { "b":1, "2": { "3": {"a":"A", "b":3, "n":[1,2]} }, "a":3 }}`,
			arrayAppend: false,
		},
		{
			src:         `{"a": ["aye"]}`,
			dst:         `{"a": ["eh"]}`,
			expected:    `{"a": ["aye"]}`,
			arrayAppend: false,
		},
		{
			src:         `{"a": ["aye"]}`,
			dst:         `{"a": ["eh"]}`,
			expected:    `{"a": ["eh", "aye"]}`,
			arrayAppend: true,
		},
	} {
		var dst map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.dst), &dst); err != nil {
			t.Fatal(err)
			continue
		}

		var src map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.src), &src); err != nil {
			t.Fatal(err)
			continue
		}

		var expected map[string]interface{}
		if err := json.Unmarshal([]byte(tuple.expected), &expected); err != nil {
			t.Fatal(err)
			continue
		}

		got := Merge(dst, src, tuple.arrayAppend)
		assert(t, expected, got)
	}
}

func assert(t *testing.T, expected, got map[string]interface{}) {
	expectedBuf, err := json.Marshal(expected)
	if err != nil {
		t.Error(err)
		return
	}
	gotBuf, err := json.Marshal(got)
	if err != nil {
		t.Error(err)
		return
	}
	if bytes.Compare(expectedBuf, gotBuf) != 0 {
		t.Errorf("expected %s, got %s", string(expectedBuf), string(gotBuf))
		return
	}
}
