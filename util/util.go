package util

func RuneToString(r rune) string {
	return string(r)
}

func StringToRune(s string, index int) rune {
	return rune(s[index])
}
