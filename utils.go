package kuroshiro

import (
	"github.com/chanyeinthaw/kuroshiro.go/checker"
	"github.com/chanyeinthaw/kuroshiro.go/converter"
)

var (
	ToHiragana = converter.ToRawHiragana
	ToKatakana = converter.ToRawKatakana
	ToRomaji   = converter.ToRawRomaji
)

var (
	IsHiragana = checker.IsHiragana
	IsKatakana = checker.IsKatakana
	IsKana     = checker.IsKana
	IsKanji    = checker.IsKanji
	IsJapanese = checker.IsJapanese
)

func HasHiragana(str string) bool {
	return checker.HasKana(str, checker.HIRAGANA)
}

func HasKatakana(str string) bool {
	return checker.HasKana(str, checker.KATAKANA)
}

func HasKana(str string) bool {
	return checker.HasKana(str, checker.KANA)
}

func HasKanji(str string) bool {
	return checker.HasKana(str, checker.KANJI)
}

func HasJapanese(str string) bool {
	return checker.HasKana(str, checker.JAPANESE)
}
