package kuroshiro

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chanyeinthaw/kuroshiro.go/analyzer"
	"github.com/chanyeinthaw/kuroshiro.go/checker"
	"github.com/chanyeinthaw/kuroshiro.go/converter"
	"golang.org/x/sync/errgroup"
)

type Kuroshiro struct {
	analyzer analyzer.Analyzer
}

type _NotionType int

const (
	_NOTION_KANJI _NotionType = 1
	_NOTION_KANA  _NotionType = 2
	_NOTION_OTHER _NotionType = 3
)

type _Notion struct {
	basic         string
	notionType    _NotionType
	notation      string
	pronunciation string
}

func New(analyzer analyzer.Analyzer) *Kuroshiro {
	return &Kuroshiro{analyzer: analyzer}
}

func (k *Kuroshiro) Convert(str string, opts *Options) (string, error) {
	if opts == nil {
		opts = NewOptions()
	}

	tokens, err := k.analyzer.Parse(str)
	if err != nil {
		return "", err
	}

	tokens = converter.PatchTokens(tokens)

	if opts.mode == NORMAL || opts.mode == SPACED {
		switch opts.to {
		case KATAKANA:
			strs := mapTokensToStr(tokens, tokenToKatakana)
			return glueStrings(strs, opts.mode), nil
		case ROMAJI:
			strs := mapTokensToStr(tokens, func(_ int, token analyzer.Token, __ []analyzer.Token) string {
				return tokenToRomaji(token, opts.romajiSystem)
			})
			return glueStrings(strs, opts.mode), nil
		case HIRAGANA:
			strs := mapTokensToStr(tokens, tokenToHiragana)
			return glueStrings(strs, opts.mode), nil
		default:
			return "", nil
		}
	}

	notations := getNotationsFromTokens(tokens)
	var result string

	switch opts.to {
	case KATAKANA:
		result = notationsToKana(opts, notations, func(notation _Notion) string {
			return converter.ToRawKatakana(notation.notation)
		})
	case ROMAJI:
		result = notationsToKana(opts, notations, func(notation _Notion) string {
			return converter.ToRawRomaji(notation.notation, opts.romajiSystem)
		})
	case HIRAGANA:
		result = notationsToKana(opts, notations, func(notation _Notion) string {
			return notation.notation
		})
	}

	return result, nil
}

func (k *Kuroshiro) ConvertMultiline(str string, opts *Options) (string, error) {
	strs := strings.Split(str, "\n")

	eg, _ := errgroup.WithContext(context.TODO())
	eg.SetLimit(1)
	results := make(chan struct {
		idx    int
		result string
	}, len(strs))

	for idx, s := range strs {
		idxx := idx
		ss := s
		eg.Go(func() error {
			if strings.Trim(ss, " ") == "" {
				results <- struct {
					idx    int
					result string
				}{
					idx:    idxx,
					result: ss,
				}
				return nil
			}

			result, err := k.Convert(ss, opts)
			if err != nil {
				return err
			}

			results <- struct {
				idx    int
				result string
			}{
				idx:    idxx,
				result: result,
			}
			return nil
		})
	}

	err := eg.Wait()
	close(results)

	if err != nil {
		return "", err
	}

	for result := range results {
		strs[result.idx] = result.result
	}

	return strings.Join(strs, "\n"), nil
}

func notationsToKana(opts *Options, notations []_Notion, convert func(_Notion) string) (result string) {
	switch opts.mode {
	case OKURIGANA:
		for _, notation := range notations {
			if notation.notionType != _NOTION_KANJI {
				result += notation.basic
			} else {
				result += fmt.Sprintf("%s%s%s%s", notation.basic, opts.delimiterStart, convert(notation), opts.delimiterEnd)
			}
		}
	case FURIGANA:
		for _, notation := range notations {
			if notation.notionType != _NOTION_KANJI {
				result += notation.basic
			} else {
				result += fmt.Sprintf("<ruby>%s<rp>%s</rp><rt>%s</rt><rp>%s</rp></ruby>", notation.basic, opts.delimiterStart, convert(notation), opts.delimiterEnd)
			}
		}
	}
	return
}

func getNotationsFromTokens(tokens []analyzer.Token) (notations []_Notion) {
	for _, token := range tokens {
		strType := checker.GetStrType(token.SurfaceForm)
		surfaceForm := []rune(token.SurfaceForm)
		switch strType {
		case checker.STR_PURE_KENJI:
			prn := token.Pronunciation
			if prn == analyzer.NUL {
				prn = token.Reading
			}

			notations = append(notations, _Notion{
				basic:         token.SurfaceForm,
				notionType:    _NOTION_KANJI,
				notation:      converter.ToRawHiragana(token.Reading),
				pronunciation: prn,
			})
		case checker.STR_KENJI_KANA:
			pattern := ""
			isLastTokenKanji := false
			var subs []rune

			for _, ch := range surfaceForm {
				if checker.IsKanji(ch) {
					if !isLastTokenKanji {
						isLastTokenKanji = true
						pattern += "(.+)"
						subs = append(subs, ch)
					} else {
						subs[len(subs)-1] = ch
					}
				} else {
					isLastTokenKanji = false
					subs = append(subs, ch)

					if checker.IsKatakana(ch) {
						pattern += converter.ToRawHiragana(string(ch))
					} else {
						pattern += string(ch)
					}
				}
			}

			reg := regexp.MustCompile(fmt.Sprintf("^%s$", pattern))
			matches := reg.FindStringSubmatch(token.Reading)

			if matches != nil {
				pickKanji := 1
				for _, sub := range subs {
					if checker.IsKanji(sub) {
						notations = append(notations, _Notion{
							basic:         string(sub),
							notionType:    _NOTION_KANJI,
							notation:      matches[pickKanji],
							pronunciation: converter.ToRawKatakana(matches[pickKanji]),
						})
						pickKanji++
					} else {
						notations = append(notations, _Notion{
							basic:         string(sub),
							notionType:    _NOTION_KANA,
							notation:      converter.ToRawHiragana(string(sub)),
							pronunciation: converter.ToRawKatakana(string(sub)),
						})
					}
				}
			} else {
				prn := token.Pronunciation
				if prn == analyzer.NUL {
					prn = token.Reading
				}

				notations = append(notations, _Notion{
					basic:         token.SurfaceForm,
					notionType:    _NOTION_KANJI,
					notation:      converter.ToRawHiragana(token.Reading),
					pronunciation: prn,
				})
			}

		case checker.STR_PURE_KANA:
			for idx, ch := range surfaceForm {
				var prn string
				if token.Reading != analyzer.NUL {
					prn = string([]rune(token.Reading)[idx])
				}

				if token.Pronunciation != analyzer.NUL {
					prn = string([]rune(token.Pronunciation)[idx])
				}

				notations = append(notations, _Notion{
					basic:         string(ch),
					notionType:    _NOTION_KANA,
					notation:      converter.ToRawHiragana(string(ch)),
					pronunciation: prn,
				})
			}
		case checker.STR_OTHER:
			for _, ch := range surfaceForm {
				notations = append(notations, _Notion{
					basic:         string(ch),
					notionType:    _NOTION_OTHER,
					notation:      string(ch),
					pronunciation: string(ch),
				})
			}
		default:
			panic("Unknown strType")
		}
	}

	return
}

func glueStrings(strs []string, mode ConvertMode) string {
	glue := ""
	if mode == SPACED {
		glue = " "
	}

	return strings.Join(strs, glue)
}

func tokenToHiragana(idx int, token analyzer.Token, tokens []analyzer.Token) string {
	if checker.HasKana(token.SurfaceForm, checker.KANJI) {
		if checker.HasKana(token.Reading, checker.KATAKANA) == false {
			tokens[idx].Reading = converter.ToRawHiragana(token.Reading)
		} else {
			tokens[idx].Reading = converter.ToRawHiragana(token.Reading)
			var tmp string
			var hpattern string
			surfaceForm := []rune(token.SurfaceForm)
			for hc := 0; hc < len(surfaceForm); hc++ {
				if checker.IsKanji(surfaceForm[hc]) {
					hpattern += "(.*)"
				} else {
					if checker.IsKatakana(surfaceForm[hc]) {
						hpattern += converter.ToRawHiragana(string(surfaceForm[hc]))
					} else {
						hpattern += string(surfaceForm[hc])
					}
				}
			}

			hreg := regexp.MustCompile(hpattern)
			hmatches := hreg.FindStringSubmatch(tokens[idx].Reading)
			if hmatches != nil {
				pickKJ := 0
				for hc1 := 0; hc1 < len(surfaceForm); hc1++ {
					if checker.IsKanji(surfaceForm[hc1]) {
						tmp += hmatches[pickKJ+1]
						pickKJ++
					} else {
						tmp += string(surfaceForm[hc1])
					}
				}

				tokens[idx].Reading = tmp
			}

		}
	} else {
		tokens[idx].Reading = token.SurfaceForm
	}
	return tokens[idx].Reading
}

func tokenToRomaji(token analyzer.Token, system converter.RomajiSystem) string {
	var preToken string
	if checker.HasKana(token.SurfaceForm, checker.JAPANESE) {
		if token.Pronunciation != analyzer.NUL {
			preToken = token.Pronunciation
		} else {
			preToken = token.Reading
		}
	} else {
		preToken = token.SurfaceForm
	}

	return converter.ToRawRomaji(preToken, system)
}

func tokenToKatakana(_ int, token analyzer.Token, __ []analyzer.Token) string {
	return token.Reading
}

func mapTokensToStr(tokens []analyzer.Token, f func(int, analyzer.Token, []analyzer.Token) string) []string {
	var result []string
	for idx, token := range tokens {
		result = append(result, f(idx, token, tokens))
	}

	return result
}
