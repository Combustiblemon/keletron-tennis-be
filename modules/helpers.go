package helpers

import (
	"strings"
)

var encoder = strings.NewReplacer("\"", "&quote")
var decoder = strings.NewReplacer("&quote", "\"")

func Encode(text string) string {
	return ""
}

func Decode(text string) string {
	return ""
}
