package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MapVar is a string map which satisfies the Value interface.
type MapVar map[string]interface{}

func (mv *MapVar) String() string {
	deref := map[string]interface{}(*mv)
	parts := make([]string, 0, len(deref))
	for name, value := range deref {
		parts = append(parts, name+"="+fmt.Sprint(value))
	}
	return strings.Join(parts, ",")
}

func (mv *MapVar) Set(value string) error {
	deref := map[string]interface{}(*mv)
	parts := strings.SplitN(value, "=", 2)
	if len(parts) == 1 {
		deref[parts[0]] = ""
	} else {
		deref[parts[0]] = parts[1]
	}
	return nil
}

// JsonVar is a JSON object which satisfies the Value interface.
type JsonVar map[string]interface{}

func (j *JsonVar) String() string {
	return fmt.Sprintf("%+v", map[string]interface{}(*j))
}

func (v *JsonVar) Set(value string) error {
	return json.Unmarshal([]byte(value), v)
}
