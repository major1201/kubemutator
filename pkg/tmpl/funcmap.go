package tmpl

import (
	"github.com/major1201/goutils"
	htmlTemplate "html/template"
	"os"
	"strings"
	textTemplate "text/template"
)

// FuncMap a common template tmpl implement
type FuncMap map[string]interface{}

// TextFuncMap convert to text template funcMap
func (fm FuncMap) TextFuncMap() textTemplate.FuncMap {
	return textTemplate.FuncMap(fm)
}

// HTMLFuncMap convert to html template funcMap
func (fm FuncMap) HTMLFuncMap() htmlTemplate.FuncMap {
	return htmlTemplate.FuncMap(fm)
}

// GetFuncMap returns the common funcMap
func GetFuncMap() FuncMap {
	return fm
}

var fm = FuncMap{
	// numbers
	"int":   goutils.ToInt,
	"intdv": goutils.ToIntDv,
	"inc":   inc,
	"add":   add,
	"sub":   sub,
	"mul":   mul,
	"div":   div,
	"mod":   mod,
	"rand":  random,

	// strings
	"title":      strings.Title,
	"replaceall": strings.ReplaceAll,
	"trim":       goutils.Trim,
	"trimleft":   goutils.TrimLeft,
	"trimright":  goutils.TrimRight,
	"upper":      strings.ToUpper,
	"lower":      strings.ToLower,
	"join":       join,
	"split":      split,
	"hasprefix":  hasPrefix,
	"hassuffix":  hasSuffix,
	"between":    goutils.Between,
	"contains":   contains,
	"indent":     indent,
	"uuid":       goutils.UUID,
	"filesize":   goutils.FileSize,
	"leftpad":    goutils.LeftPad,
	"rightpad":   goutils.RightPad,

	// bool
	"bool":   goutils.ToBool,
	"booldv": goutils.ToBoolDv,

	// encoding
	"base64en":   encodeBase64,
	"base64de":   decodeBase64,
	"md5":        encodeMd5,
	"sha1":       encodeSha1,
	"sha224":     encodeSha224,
	"sha256":     encodeSha256,
	"sha512":     encodeSha512,
	"json":       toJSON,
	"prettyjson": toPrettyJSON,
	"yaml":       toYAML,

	// system
	"debug": debug,
	"env":   os.Getenv,
	"idx":   index,
}
