package pangoutils

import "strings"

func Chomp(s string) string {
	return strings.TrimRight(s, "\n\t ")
}

func Chompb(s []byte) string {
	return Chomp(string(s))
}