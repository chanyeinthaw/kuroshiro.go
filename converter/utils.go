package converter

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func substr_end(input string, start int, end int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if end > len(asRunes) {
		end = len(asRunes)
	}

	return string(asRunes[start:end])
}
