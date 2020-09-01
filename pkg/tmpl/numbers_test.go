package tmpl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNumbers(t *testing.T) {
	tmpl := `
{{ "77" | int }}
{{ "33" | int | inc }}
{{ add 7 3.45 }}
{{ sub 7 3 }}
{{ mul 7 3 }}
{{ div 7 3 | printf "%.2f" }}
{{ mod 7 3 }}
`

	expect := `
77
34
10.45
4
21
2.33
1
`

	ta := assert.New(t)

	result, err := ExecuteTextTemplate(tmpl, nil)
	ta.NoError(err)
	if err != nil {
		return
	}

	ta.Equal(expect, result)
}
