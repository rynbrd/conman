package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MapVar is a string map which satisfies the Value interface.
type MapVar struct {
	Context map[string]interface{}
}

func (v *MapVar) String() string {
	parts := make([]string, 0, len(v.Context))
	for name, value := range v.Context {
		parts = append(parts, name+"="+fmt.Sprint(value))
	}
	return strings.Join(parts, ",")
}

func (v *MapVar) Set(value string) error {
	if v.Context == nil {
		v.Context = map[string]interface{}{}
	}
	parts := strings.SplitN(value, "=", 2)
	if len(parts) == 1 {
		v.Context[parts[0]] = ""
	} else {
		v.Context[parts[0]] = parts[1]
	}
	return nil
}

// JsonVar is a JSON object which satisfies the Value interface.
type JsonVar struct {
	Context map[string]interface{}
}

func (j *JsonVar) String() string {
	return fmt.Sprintf("%+v", j.Context)
}

func (v *JsonVar) Set(value string) error {
	if v.Context == nil {
		v.Context = map[string]interface{}{}
	}
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(value), &m)
	if err == nil {
		v.Context = Merge(v.Context, m, true)
	}
	return err
}
