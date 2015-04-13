package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var TemplateFuncs template.FuncMap = template.FuncMap{
	"replace": strings.Replace,
	"join":    strings.Join,
	"split":   strings.Split,
	"title":   strings.Title,
	"upper":   strings.ToUpper,
	"lower":   strings.ToLower,
	"json":    JSON,
}

// Template represents a single template to be rendered by ConMan.
type Template struct {
	Src  string
	Dest string
}

// Render the template.
func (t *Template) Render(context map[string]interface{}) error {
	wrapError := func(err error) error {
		return fmt.Errorf("%s: %s", t.Src, err)
	}

	// create the destination directory
	dir := filepath.Dir(t.Dest)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("%s: %s", t.Dest, err)
	}

	// render the template
	name := filepath.Base(t.Src)
	if dest, err := os.Create(t.Dest); err == nil {
		if tpl, err := template.New(name).Funcs(TemplateFuncs).ParseFiles(t.Src); err == nil {
			return tpl.Execute(dest, context)
		} else {
			return wrapError(err)
		}
	} else {
		return wrapError(err)
	}
}

// RenderString takes a template string and renders it using the provided context.
func RenderString(value string, context map[string]interface{}) (string, error) {
	wrapError := func(err error) error {
		return fmt.Errorf("'%s': %s", value, err)
	}

	if tpl, err := template.New("string").Funcs(TemplateFuncs).Parse(value); err == nil {
		buf := &bytes.Buffer{}
		if err := tpl.Execute(buf, context); err == nil {
			return buf.String(), nil
		} else {
			return "", wrapError(err)
		}
	} else {
		return "", wrapError(err)
	}
}

// JSON parses the JSON value `item` and returns the result of indexing the
// resulting value by the following arguments.
func JSON(item interface{}, indices ...interface{}) (interface{}, error) {
	var data []byte
	switch item.(type) {
	case []byte:
		data = item.([]byte)
	case string:
		data = []byte(item.(string))
	default:
		return nil, errors.New("invalid input type")
	}

	values := map[string]interface{}{}
	if err := json.Unmarshal(data, values); err != nil {
		return nil, err
	}
	return Index(values, indices...)
}

// Index returns the result of indexing its first argument by the following
// arguments.  Thus "index x 1 2 3" is, in Go syntax, x[1][2][3]. Each
// indexed item must be a map, slice, or array.
//
// Borrowed from the standard library.
func Index(item interface{}, indices ...interface{}) (interface{}, error) {
	indirect := func(v reflect.Value) (rv reflect.Value, isNil bool) {
		for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
			if v.IsNil() {
				return v, true
			}
			if v.Kind() == reflect.Interface && v.NumMethod() > 0 {
				break
			}
		}
		return v, false
	}

	v := reflect.ValueOf(item)
	for _, i := range indices {
		index := reflect.ValueOf(i)
		var isNil bool
		if v, isNil = indirect(v); isNil {
			return nil, fmt.Errorf("index of nil pointer")
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			var x int64
			switch index.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x = index.Int()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				x = int64(index.Uint())
			default:
				return nil, fmt.Errorf("cannot index slice/array with type %s", index.Type())
			}
			if x < 0 || x >= int64(v.Len()) {
				return nil, fmt.Errorf("index out of range: %d", x)
			}
			v = v.Index(int(x))
		case reflect.Map:
			if !index.IsValid() {
				index = reflect.Zero(v.Type().Key())
			}
			if !index.Type().AssignableTo(v.Type().Key()) {
				return nil, fmt.Errorf("%s is not index type for %s", index.Type(), v.Type())
			}
			if x := v.MapIndex(index); x.IsValid() {
				v = x
			} else {
				v = reflect.Zero(v.Type().Elem())
			}
		default:
			return nil, fmt.Errorf("can't index item of type %s", v.Type())
		}
	}
	return v.Interface(), nil
}
