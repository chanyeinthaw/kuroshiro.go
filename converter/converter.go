package converter

import (
	"regexp"
	"unicode/utf8"
)

const (
	KATAKANA_HIRAGANA_SHIFT = '\u3041' - '\u30a1'
	HIRAGANA_KATAKANA_SHIFT = '\u30a1' - '\u3041'
)

func ToRawHiragana(str string) string {
	runes := make([]rune, len(str))

	for _, ch := range str {
		if ch >= '\u30a0' && ch <= '\u30f7' {
			runes = append(runes, ch+KATAKANA_HIRAGANA_SHIFT)
		} else {
			runes = append(runes, ch)
		}
	}

	return string(runes)
}

func ToRawKatakana(str string) string {
	runes := make([]rune, len(str))

	for _, ch := range str {
		if ch >= '\u3040' && ch <= '\u3097' {
			runes = append(runes, ch+HIRAGANA_KATAKANA_SHIFT)
		} else {
			runes = append(runes, ch)
		}
	}

	return string(runes)
}

func ToRawRomaji(str string, system RomajiSystem) string {
	regTsu := regexp.MustCompile(`(?m)(っ|ッ)([bcdfghijklmnopqrstuvwyz])`)
	regXTsu := regexp.MustCompile(`(?m)っ|ッ`)

	romajiSystem := romajiSystems[system]

	pnt := 0
	result := ""

	if system == ROMAJI_PASSPORT {
		re := regexp.MustCompile(`(?m)ー`)
		str = re.ReplaceAllString(str, "ｰ")
	}

	if system == ROMAJI_NIPPON || system == ROMAJI_HEPBURN {
		regHatu := regexp.MustCompile(`(ん|ン)([あいうえおアィウエオぁぃぅぇぉァィゥェォやゆよヤユヨゃゅょャュョ])`)

		var indices []int
		for _, matches := range regHatu.FindAllStringSubmatchIndex(str, -1) {
			runeIndex := utf8.RuneCountInString(str[:matches[2]])
			indices = append(indices, runeIndex+1)
		}

		if len(indices) != 0 {
			mStr := ""
			for i, idx := range indices {
				if i == 0 {
					mStr += substr_end(str, 0, idx) + "'"
				} else {
					mStr += substr_end(str, indices[i-1], idx) + "'"
				}
			}
			mStr += substr_end(str, indices[len(indices)-1], len(str))
			str = mStr
		}
	}

	// [ALL] kana to roman chars
	max := utf8.RuneCountInString(str)
	for pnt <= max {
		ch := substr(str, pnt, 2)
		if r, ok := romajiSystem[ch]; ok {
			result += r
			pnt += 2
		} else {
			ch = substr(str, pnt, 1)
			if r, ok := romajiSystem[ch]; ok {
				result += r
			} else {
				result += ch
			}
			pnt += 1
		}
	}
	result = regTsu.ReplaceAllString(result, "$2$2")

	// [PASSPORT|HEPBURN] 子音を重ねて特殊表記
	if system == ROMAJI_PASSPORT || system == ROMAJI_HEPBURN {
		re := regexp.MustCompile(`(?m)cc`)
		result = re.ReplaceAllString(result, "tc")
	}

	result = regXTsu.ReplaceAllString(result, "tsu")

	if system == ROMAJI_PASSPORT || system == ROMAJI_HEPBURN {
		re := regexp.MustCompile(`(?m)nm`)
		result = re.ReplaceAllString(result, "mm")

		re = regexp.MustCompile(`(?m)nb`)
		result = re.ReplaceAllString(result, "mb")

		re = regexp.MustCompile(`(?m)np`)
		result = re.ReplaceAllString(result, "mp")
	}

	if system == ROMAJI_NIPPON {
		re := regexp.MustCompile(`(?m)aー`)
		result = re.ReplaceAllString(result, "â")

		re = regexp.MustCompile(`(?m)iー`)
		result = re.ReplaceAllString(result, "î")

		re = regexp.MustCompile(`(?m)uー`)
		result = re.ReplaceAllString(result, "û")

		re = regexp.MustCompile(`(?m)eー`)
		result = re.ReplaceAllString(result, "ê")

		re = regexp.MustCompile(`(?m)oー`)
		result = re.ReplaceAllString(result, "ô")
	}

	if system == ROMAJI_HEPBURN {
		re := regexp.MustCompile(`(?m)aー`)
		result = re.ReplaceAllString(result, "ā")

		re = regexp.MustCompile(`(?m)iー`)
		result = re.ReplaceAllString(result, "ī")

		re = regexp.MustCompile(`(?m)uー`)
		result = re.ReplaceAllString(result, "ū")

		re = regexp.MustCompile(`(?m)eー`)
		result = re.ReplaceAllString(result, "ē")

		re = regexp.MustCompile(`(?m)oー`)
		result = re.ReplaceAllString(result, "ō")
	}

	return result
}
