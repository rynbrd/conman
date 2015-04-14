package main

// Context managed a map of context data for use in template rendering.
type Context map[string]interface{}

// Update the context with additional content. This works by merging the new
// tree of values with the existing context tree.
func (c *Context) Update(values map[string]interface{}) {
	Merge(*c, values)
}

// Map returns the context as a map suitable for template rendering.
func (c *Context) Map() map[string]interface{} {
	return map[string]interface{}(*c)
}
