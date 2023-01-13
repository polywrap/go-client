package msgpack

import "strings"

func Capitalize(v string) string {
	raw := []rune(v)
	return strings.ToUpper(string(raw[0])) + string(raw[1:])
}

func UnCapitalize(v string) string {
	raw := []rune(v)
	return strings.ToLower(string(raw[0])) + string(raw[1:])
}
