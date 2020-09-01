package tmpl

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func decodeBase64(s string) string {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func encodeMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func encodeSha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

func encodeSha224(s string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(s)))
}

func encodeSha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func encodeSha512(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func toPrettyJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func toYAML(v interface{}) string {
	b, _ := yaml.Marshal(v)
	return string(b)
}
