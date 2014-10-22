package main

import (
	"bytes"
	"testing"

	assert "github.com/pilu/miniassert"
)

func TestParseConfig(t *testing.T) {
	content := `javascripts:
  application:
    - file-1.js
    - file-2.js
  admin:
    - file-3.js
    - file-4.js
`
	r := bytes.NewBufferString(content)

	c, err := parseConfig(r)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(c.Javascripts))

	assert.Type(t, "[]string", c.Javascripts["application"])
	assert.Equal(t, 2, len(c.Javascripts["application"]))
	assert.Equal(t, "file-1.js", c.Javascripts["application"][0])
	assert.Equal(t, "file-2.js", c.Javascripts["application"][1])

	assert.Type(t, "[]string", c.Javascripts["admin"])
	assert.Equal(t, 2, len(c.Javascripts["admin"]))
	assert.Equal(t, "file-3.js", c.Javascripts["admin"][0])
	assert.Equal(t, "file-4.js", c.Javascripts["admin"][1])
}

func TestParseConfig_UnmarshalError(t *testing.T) {
	content := `Unmarshal\nError`
	r := bytes.NewBufferString(content)

	_, err := parseConfig(r)
	assert.NotNil(t, err)
}
