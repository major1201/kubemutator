package tmpl

import (
	"github.com/major1201/goutils"
	"strings"
)

func join(sep string, elems []string) string {
	return strings.Join(elems, sep)
}

func split(sep, s string) []string {
	return strings.Split(s, sep)
}

func hasPrefix(prefix, s string) bool {
	return strings.HasPrefix(s, prefix)
}

func hasSuffix(suffix, s string) bool {
	return strings.HasSuffix(s, suffix)
}

func contains(substr, s string) bool {
	return strings.Contains(s, substr)
}

func indent(indent int, s string) string {
	return goutils.Indent(s, strings.Repeat(" ", indent))
}
