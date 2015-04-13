package main

import (
	"strings"
)

// Environ is a collection of environment variables. It does not maintain the
// order of its values.
type Environ map[string]string

// Update values in the environment. Existing values are overwritten.
func (env *Environ) Update(environ map[string]string) {
	for name, value := range environ {
		(map[string]string(*env))[name] = value
	}
}

// Load values into the environment. The elements in the `environ` param must
// be formatted as `NAME=VALUE` as they are split on the first occurence of
// `=`. Duplicate names within the list are ignored. Existing names in `env`
// are overwritten.
func (env *Environ) Load(environ []string) {
	envMap := make(map[string]string, len(environ))
	for _, envVar := range environ {
		name, value := ParseEnvVar(envVar)
		if _, ok := envMap[name]; !ok {
			envMap[name] = value
		}
	}
	env.Update(envMap)
}

// Values returns the Environ as a array of NAME=VALUE formatted strings.
func (env *Environ) Values() []string {
	values := make([]string, 0, len(*env))
	for name, value := range map[string]string(*env) {
		values = append(values, EncodeEnvVar(name, value))
	}
	return values
}

// Context returns the environment as a map suitable for use as template context.
func (env *Environ) Context() map[string]interface{} {
	envMap := map[string]string(*env)
	context := make(map[string]interface{}, len(envMap))
	for name, value := range envMap {
		context[name] = value
	}
	return context
}

// ParseEnvVar parses a single environment variable.
func ParseEnvVar(envVar string) (string, string) {
	parts := strings.SplitN(envVar, "=", 2)
	if len(parts) == 1 {
		return parts[0], ""
	} else {
		return parts[0], parts[1]
	}
}

// EncodeEnvVar returns a
func EncodeEnvVar(name, value string) string {
	return name + "=" + value
}
