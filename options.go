package kuroshiro

import "github.com/chanyeinthaw/kuroshiro.go/converter"

type ConvertTo string
type ConvertMode string

const (
	HIRAGANA ConvertTo = "hiragana"
	KATAKANA ConvertTo = "katakana"
	ROMAJI   ConvertTo = "romaji"
)

const (
	NORMAL    ConvertMode = "normal"
	SPACED    ConvertMode = "spaced"
	OKURIGANA ConvertMode = "okurigana"
	FURIGANA  ConvertMode = "furigana"
)

type RomajiSystem string

const (
	NIPPON   = converter.ROMAJI_NIPPON
	PASSPORT = converter.ROMAJI_PASSPORT
	HEPBURN  = converter.ROMAJI_HEPBURN
)

type Options struct {
	to             ConvertTo
	mode           ConvertMode
	romajiSystem   converter.RomajiSystem
	delimiterStart string
	delimiterEnd   string
}

func NewOptions() *Options {
	return &Options{
		to:             HIRAGANA,
		mode:           NORMAL,
		romajiSystem:   HEPBURN,
		delimiterStart: "(",
		delimiterEnd:   ")",
	}
}

func (o *Options) ConvertTo(to ConvertTo) *Options {
	o.to = to
	return o
}

func (o *Options) SetMode(mode ConvertMode) *Options {
	o.mode = mode
	return o
}

func (o *Options) SetRomajiSystem(system converter.RomajiSystem) *Options {
	o.romajiSystem = system
	return o
}

func (o *Options) SetDelimiter(start, end string) *Options {
	o.delimiterStart = start
	o.delimiterEnd = end
	return o
}
