package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	assert "github.com/pilu/miniassert"
)

func TestJSMinify(t *testing.T) {
	content := `function() {
		return 1
	}`
	r := bytes.NewBufferString(content)
	b, err := jsMinify(r)

	assert.Nil(t, err)
	assert.Equal(t, "\nfunction(){return 1}", string(b))
}

func TestNewMinifier(t *testing.T) {
	r1 := bytes.NewBufferString("function a() {\n return 'a'}")
	r2 := bytes.NewBufferString("function b() {\n return 'b'}")
	sources := []io.Reader{r1, r2}

	out := bytes.NewBuffer([]byte{})
	m := newMinifier(jsMinify, sources, out)

	assert.Equal(t, reflect.ValueOf(jsMinify).Pointer(), reflect.ValueOf(m.mf).Pointer())

	err := m.minify()

	assert.Nil(t, err)
	assert.Equal(t, "\nfunction a(){return'a'}\nfunction b(){return'b'}", string(out.Bytes()))
}
