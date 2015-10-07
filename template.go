package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var TemplateFuncs template.FuncMap = template.FuncMap{
	"replace":     strings.Replace,
	"join":        strings.Join,
	"split":       strings.Split,
	"title":       strings.Title,
	"upper":       strings.ToUpper,
	"lower":       strings.ToLower,
	"trim":        strings.Trim,
	"trimSpace":   strings.TrimSpace,
	"json":        JSON,
	"addrHost":    AddrHost,
	"addrPort":    AddrPort,
	"urlScheme":   URLScheme,
	"urlUsername": URLUsername,
	"urlPassword": URLPassword,
	"urlHost":     URLHost,
	"urlPath":     URLPath,
	"urlRawQuery": URLRawQuery,
	"urlQuery":    URLQuery,
	"urlFragment": URLFragment,
}

// Template represents a single template to be rendered by ConMan.
type Template struct {
	Src string
	Dst string
}

// Render the template.
func (t *Template) Render(context map[string]interface{}) error {
	wrapError := func(err error) error {
		return fmt.Errorf("%s: %s", t.Src, err)
	}

	// create the destination directory
	dir := filepath.Dir(t.Dst)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("%s: %s", t.Dst, err)
	}

	// render the template
	name := filepath.Base(t.Src)
	if dest, err := os.Create(t.Dst); err == nil {
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
