package tmpl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKubernetes(t *testing.T) {
	tmpl := `
{{ mulquantity 2 "500m" }}
{{ 5.3 | mulquantity 2 }}
{{ "5.3" | mulquantity 2 }}
{{ divquantity 2 "500m" }}
{{ 5.3 | divquantity 2 }}
`

	expect := `
1
10600m
10600m
250m
2650m
`

	ta := assert.New(t)

	result, err := ExecuteTextTemplate(tmpl, nil)
	ta.NoError(err)
	if err != nil {
		return
	}

	ta.Equal(expect, result)
}
