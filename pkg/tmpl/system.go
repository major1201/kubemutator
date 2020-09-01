package tmpl

import (
	"fmt"
	"strings"
)

func debug(any ...interface{}) string {
	s := make([]string, len(any))
	for i, a := range any {
		s[i] = fmt.Sprintf("%[1]T %[1]v", a)
	}
	return strings.Join(s, " ")
}

func index(i int, a interface{}) interface{} {
	if a == nil {
		return nil
	}
	switch a := a.(type) {
	case []string:
		if i < 0 || i >= len(a) {
			return -1
		}
		return a[i]
	case []int64:
		if i < 0 || i >= len(a) {
			return -1
		}
		return a[i]
	case string:
		if i < 0 || i >= len(a) {
			return -1
		}
		return fmt.Sprintf("%c", a[i])
	}
	return a
}
