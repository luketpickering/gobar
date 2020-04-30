package pangoutils

func chomp(s string) string {
	return strings.TrimRight(s, "\n\t ")
}

func chompb(s []byte) string {
	return chomp(string(s))
}