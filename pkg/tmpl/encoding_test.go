package tmpl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestData struct {
	Int    int     `json:"int"`
	Float  float64 `json:"float"`
	String string  `json:"string"`
}

func TestEncoding(t *testing.T) {
	tmpl := `
{{ "hello" | base64en }}
{{ "aGVsbG8=" | base64de }}
{{ "hello" | md5 }}
{{ "hello" | sha1 }}
{{ "hello" | sha224 }}
{{ "hello" | sha256 }}
{{ "hello" | sha512 }}
{{ . | json }}
{{ . | prettyjson }}
{{ . | yaml }}
`

	expect := `
aGVsbG8=
hello
5d41402abc4b2a76b9719d911017c592
aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
ea09ae9cc6768c50fcee903ed054556e5bfc8347907f12598aa24193
2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043
{"int":33,"float":1.23,"string":"hello world"}
{
  "int": 33,
  "float": 1.23,
  "string": "hello world"
}
int: 33
float: 1.23
string: hello world

`

	ta := assert.New(t)

	data := TestData{
		Int:    33,
		Float:  1.23,
		String: "hello world",
	}

	result, err := ExecuteTextTemplate(tmpl, data)
	ta.NoError(err)
	if err != nil {
		return
	}

	ta.Equal(expect, result)
}
