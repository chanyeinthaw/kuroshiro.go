package checker

func IsHiragana(ch rune) bool {
	return ch >= '\u3040' && ch <= '\u309f'
}

func IsKatakana(ch rune) bool {
	return ch >= '\u30a0' && ch <= '\u30ff'
}

func IsKana(ch rune) bool {
	return IsHiragana(ch) || IsKatakana(ch)
}

func IsKanji(ch rune) bool {
	return (ch >= '\u4e00' && ch <= '\u9fcf') || (ch >= '\uf900' && ch <= '\ufaff') || (ch >= '\u3400' && ch <= '\u4dbf')
}

func IsJapanese(ch rune) bool {
	return IsKana(ch) || IsKanji(ch)
}

type KanaType int

const (
	HIRAGANA KanaType = iota
	KATAKANA
	KANA
	KANJI
	JAPANESE
)

var checkerMap = map[KanaType]func(rune) bool{
	HIRAGANA: IsHiragana,
	KATAKANA: IsKatakana,
	KANA:     IsKana,
	KANJI:    IsKanji,
	JAPANESE: IsJapanese,
}

func HasKana(str string, t KanaType) bool {
	checker := checkerMap[t]

	for _, ch := range str {
		if checker(ch) {
			return true
		}
	}
	return false
}

type _StrType int

const (
	STR_PURE_KENJI _StrType = iota
	STR_KENJI_KANA
	STR_PURE_KANA
	STR_OTHER
)

func GetStrType(str string) _StrType {
	hasKJ := false
	hasHK := false

	for _, ch := range str {
		if IsKanji(ch) {
			hasKJ = true
		} else if IsKana(ch) {
			hasHK = true
		}
	}

	if hasKJ && hasHK {
		return STR_KENJI_KANA
	}

	if hasKJ {
		return STR_PURE_KENJI
	}

	if hasHK {
		return STR_PURE_KANA
	}

	return STR_OTHER
}
