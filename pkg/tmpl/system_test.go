package tmpl

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSystem(t *testing.T) {
	tmpl := `
{{ env "mytest" }}
{{ "aa,bb,cc,dd" | split "," | idx 2 }}
{{ "aa,bb,cc,dd" | split "," | len }}
`

	expect := `
ttttt
cc
4
`

	ta := assert.New(t)

	ta.NoError(os.Setenv("mytest", "ttttt"))

	result, err := ExecuteTextTemplate(tmpl, nil)
	ta.NoError(err)
	if err != nil {
		return
	}

	ta.Equal(expect, result)
}
