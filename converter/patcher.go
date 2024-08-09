package converter

import (
	"fmt"
	"unicode/utf8"

	"github.com/chanyeinthaw/kuroshiro.go/analyzer"
	"github.com/chanyeinthaw/kuroshiro.go/checker"
)

func checkEveryRune(str string, callback func(rune) bool) bool {
	for _, v := range str {
		if !callback(v) {
			return false
		}
	}
	return true
}

func getLastRune(str string) (int, rune) {
	runes := []rune(str)
	idx := len(runes) - 1
	return idx, runes[idx]
}

func PatchTokens(tokens []analyzer.Token) []analyzer.Token {
	// patch for token structure
	for cr, token := range tokens {
		if checker.HasKana(token.SurfaceForm, checker.JAPANESE) {
			if token.Reading == analyzer.NUL {
				if checkEveryRune(token.SurfaceForm, checker.IsKana) {
					tokens[cr].Reading = ToRawKatakana(token.SurfaceForm)
				} else {
					tokens[cr].Reading = token.SurfaceForm
				}
			} else if !checker.HasKana(token.Reading, checker.HIRAGANA) {
				tokens[cr].Reading = ToRawKatakana(token.Reading)
			}
		} else {
			tokens[cr].Reading = token.SurfaceForm
		}
	}

	// patch for 助動詞"う" after 動詞
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Pos == "助動詞" && (tokens[i].SurfaceForm == "う" || tokens[i].SurfaceForm == "ウ") {
			if i-1 >= 0 && tokens[i-1].Pos == "動詞" {
				tokens[i-1].SurfaceForm += "う"
				if tokens[i-1].Pronunciation != analyzer.NUL {
					tokens[i-1].Pronunciation += "ー"
				} else {
					tokens[i-1].Pronunciation = fmt.Sprintf("%sー", tokens[i-1].Reading)
				}
				tokens[i-1].Reading += "ウ"
				tokens = append(tokens[:i], tokens[i+1:]...)
				i--
			}
		}
	}

	// patch for "っ" at the tail of 動詞、形容詞
	for j := 0; j < len(tokens); j++ {
		sFormLen := utf8.RuneCountInString(tokens[j].SurfaceForm)
		_, lastRune := getLastRune(tokens[j].SurfaceForm)

		if (tokens[j].Pos == "動詞" || tokens[j].Pos == "形容詞") && sFormLen > 1 && (string(lastRune) == "っ" || string(lastRune) == "ッ") {
			if j+1 < len(tokens) {
				tokens[j].SurfaceForm += tokens[j+1].SurfaceForm
				if tokens[j].Pronunciation != analyzer.NUL {
					tokens[j].Pronunciation += tokens[j+1].Pronunciation
				} else {
					tokens[j].Pronunciation = fmt.Sprintf("%s%s", tokens[j].Reading, tokens[j+1].Reading)
				}
				tokens[j].Reading += tokens[j+1].Reading
				tokens = append(tokens[:j], tokens[j+1:]...)
				j--
			}
		}
	}

	return tokens
}
