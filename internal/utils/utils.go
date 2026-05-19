package utils

func FirstLetterToUpper(s string) string {
	if len(s) != 0 && s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}

	return s
}
