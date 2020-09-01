package tmpl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStrings(t *testing.T) {
	tmpl := `
{{ title "hello world!" }}
{{ replaceall "abababca" "ba" "xy" }}
{{ "  hello " | trim }}.txt
{{ "  hello " | trimleft }}.txt
{{ "  hello " | trimright }}.txt
{{ "Hello World!" | upper }}
{{ "Hello World!" | lower }}
{{ "aa,bb,cc,dd" | split "," | join "|" }}
{{ "important" | hasprefix "im" }}
{{ "important" | hasprefix "in" }}
{{ "important" | hassuffix "tant" }}
{{ "important" | hassuffix "dent" }}
{{ between "mytool<p>this is a par</p>fefsfea" "<p>" "</p>" }}
{{ "important" | contains "por" }}
{{ . | prettyjson | indent 2 }}
{{ filesize 115797 }}
{{ leftpad "113" "0" 5 }}
{{ rightpad "113" "0" 5 }}
`

	expect := `
Hello World!
axyxybca
hello.txt
hello .txt
  hello.txt
HELLO WORLD!
hello world!
aa|bb|cc|dd
true
false
true
false
this is a par
true
  {
    "int": 33,
    "float": 1.23,
    "string": "hello world"
  }
113 KB
00113
11300
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
