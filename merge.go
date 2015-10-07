package main

import (
	"fmt"
	"reflect"
)

var (
	MaxDepth = 32
)

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
// Arrays are merged by appending src to dst.
func Merge(dst, src map[string]interface{}) map[string]interface{} {
	return merge(dst, src, 0)
}

func merge(dst, src map[string]interface{}, depth int) map[string]interface{} {
	if depth > MaxDepth {
		panic("too deep!")
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := mapify(srcVal)
			dstMap, dstMapOk := mapify(dstVal)
			srcArr, srcArrOk := arrayify(srcVal)
			dstArr, dstArrOk := arrayify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = merge(dstMap, srcMap, depth+1)
			} else if srcArrOk && dstArrOk {
				srcVal = append(dstArr, srcArr...)
			} else {
			}
		}
		dst[key] = srcVal
	}
	return dst
}

func mapify(i interface{}) (map[string]interface{}, bool) {
	value := reflect.ValueOf(i)
	m := map[string]interface{}{}
	ok := false

	if value.Kind() == reflect.Map {
		m = make(map[string]interface{}, value.Len())
		for _, k := range value.MapKeys() {
			m[fmt.Sprintf("%s", k.Interface())] = value.MapIndex(k).Interface()
		}
		ok = true
	}
	return m, ok
}

func arrayify(i interface{}) ([]interface{}, bool) {
	value := reflect.ValueOf(i)
	kind := value.Kind()
	a := []interface{}{}
	ok := false

	if kind == reflect.Array || kind == reflect.Slice {
		a = make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			a[i] = value.Index(i).Interface()
		}
		ok = true
	}
	return a, ok
}
