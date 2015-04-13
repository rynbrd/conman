package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestTemplate(t *testing.T) {
	tmp, err := ioutil.TempDir("", "conman_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(tmp)
	}()

	src := path.Join(tmp, "src")
	dest := path.Join(tmp, "dest")

	// valid template
	if f, err := os.Create(src); err != nil {
		t.Fatal(err)
	} else if _, err := f.Write([]byte(`{{ .greeting }}, {{ .subject }}!`)); err != nil {
		t.Fatal(err)
	} else {
		f.Close()
	}

	want := "Hello, world!"
	context := map[string]interface{}{"greeting": "Hello", "subject": "world"}
	tpl := &Template{src, dest}
	if err := tpl.Render(context); err != nil {
		t.Error(err)
	}
	if have, err := ioutil.ReadFile(dest); err != nil {
		t.Error(err)
	} else if string(have) != want {
		t.Errorf("'%s' != '%s'", have, want)
	}

	// invalid template
	if f, err := os.Create(src); err != nil {
		t.Fatal(err)
	} else if _, err := f.Write([]byte(`{{ Oops!`)); err != nil {
		t.Fatal(err)
	} else {
		f.Close()
	}

	tpl = &Template{src, dest}
	if err := tpl.Render(context); err == nil {
		t.Error("no error")
	}
}

func TestRenderString(t *testing.T) {
	// no substitutions
	context := map[string]interface{}{}
	tpl := "Hello, world!"
	want := "Hello, world!"
	if have, err := RenderString(tpl, context); err != nil {
		t.Error(err)
	} else if have != want {
		t.Errorf("'%s' != '%s'", have, want)
	}

	// a substitution
	context = map[string]interface{}{"greeting": "Hello", "subject": "world"}
	tpl = "{{ .greeting }}, {{ .subject }}!"
	want = "Hello, world!"
	if have, err := RenderString(tpl, context); err != nil {
		t.Error(err)
	} else if have != want {
		t.Errorf("'%s' != '%s'", have, want)
	}

	// invalid template
	context = map[string]interface{}{}
	tpl = "{{ Oops!"
	if _, err := RenderString(tpl, context); err == nil {
		t.Error("no error")
	} else {
		t.Log(err)
	}
}
