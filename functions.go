package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// AddrHost takes an addr string of the form host:port and returns the host
// part.
func AddrHost(addr string) string {
	parts := strings.SplitN(addr, ":", 2)
	return parts[0]
}

// AddrPort takes an addr string of the form host:port and returns the port
// part. If no port exists an empty string is returned.
func AddrPort(addr string) string {
	parts := strings.SplitN(addr, ":", 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// URLScheme parses the URL and returns the scheme. An empty string is returned
// if the URL is not parseable.
func URLScheme(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		return parsed.Scheme
	}
	return ""
}

// URLUsername parses the URL and returns the username. An empty string is
// returned if the URL is not parseable.
func URLUsername(u string) string {
	if parsed, err := url.Parse(u); err == nil && parsed.User != nil {
		return parsed.User.Username()
	}
	return ""
}

// URLPassword parses the URL and returns the password. An empty string is
// returned if the URL is not parseable.
func URLPassword(u string) string {
	if parsed, err := url.Parse(u); err == nil && parsed.User != nil {
		if pass, ok := parsed.User.Password(); ok {
			return pass
		}
	}
	return ""
}

// URLHost parses the URL and returns the host. An empty string is returned
// if the URL is not parseable.
func URLHost(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		return parsed.Host
	}
	return ""
}

// URLPath parses the URL and returns the path. An empty string is returned
// if the URL is not parseable.
func URLPath(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		return parsed.Path
	}
	return ""
}

// URLRawQuery parses the URL and returns the full query string. An empty
// string is returned if the URL is not parseable.
func URLRawQuery(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		return parsed.RawQuery
	}
	return ""
}

// URLQuery parses the URL and query string and returns the first value of the
// named query variable. An empty string is returned if parsing fails or the
// query value does not exist.
func URLQuery(u, name string) string {
	if parsed, err := url.Parse(u); err == nil {
		if query, err := url.ParseQuery(parsed.RawQuery); err == nil {
			return query.Get(name)
		}
	}
	return ""
}

// URLFragment parses the URL and returns the fragment. An empty string is
// returned if the URL is not parseable.
func URLFragment(u string) string {
	if parsed, err := url.Parse(u); err == nil {
		return parsed.Fragment
	}
	return ""
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
	if err := json.Unmarshal(data, &values); err != nil {
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
