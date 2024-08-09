
![kuroshiro](https://kuroshiro.org/kuroshiro.png)

# kuroshiro.go

[![License](https://img.shields.io/github/license/lassjs/lass.svg)](LICENSE)

kuroshiro.go is a Golang port of [kuroshiro](https://github.com/hexenq/kuroshiro), a Japanese language library for converting Japanese sentence to Hiragana, Katakana or Romaji with furigana and okurigana modes supported.

## Feature
- Japanese Sentence => Hiragana, Katakana or Romaji
- Furigana and okurigana supported
- Multiple romanization systems supported
- Useful Japanese utils
    
## Prerequisites
kuroshiro.go uses [mecab](https://taku910.github.io/mecab) internally. For install instructions of mecab, you could check the official website of mecab from here.

You need to tell Go where MeCab has been installed.
```sh
$ export CGO_LDFLAGS="-L/path/to/lib -lmecab -lstdc++"
$ export CGO_CFLAGS="-I/path/to/include"
```

If you installed mecab-config, execute following comands.
```sh
$ export CGO_LDFLAGS="`mecab-config --libs`"
$ export CGO_CFLAGS="-I`mecab-config --inc-dir`"
```

## Usage
Install with go get
```sh
$ go get github.com/chanyeinthaw/kuroshiro.go
```

```js
import (
	"fmt"

	"github.com/chanyeinthaw/kuroshiro.go"
	"github.com/chanyeinthaw/kuroshiro.go/analyzer"
)

const INPUT = "感じ取れたら手を繋ごう、重なるのは人生のライン and レミリア最高！"
func main() {
	analyzer, err := analyzer.NewMecab()
	defer analyzer.Destroy()
	if err != nil {
		panic(err)
	}

	ks := kuroshiro.New(analyzer)
	opts := kuroshiro.NewOptions().ConvertTo(kuroshiro.HIRAGANA).SetMode(kuroshiro.SPACED)
	result, err := ks.Convert(INPUT, opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
```
 
## Romanization System
kuroshiro supports three kinds of romanization systems.

`nippon`: Nippon-shiki romanization. Refer to [ISO 3602 Strict](http://www.age.ne.jp/x/nrs/iso3602/iso3602.html).

`passport`: Passport-shiki romanization. Refer to [Japanese romanization table](https://www.ezairyu.mofa.go.jp/passport/hebon.html) published by Ministry of Foreign Affairs of Japan.

`hepburn`: Hepburn romanization. Refer to [BS 4812 : 1972](https://archive.is/PiJ4).

There is a useful [webpage](http://jgrammar.life.coocan.jp/ja/data/rohmaji2.htm) for you to check the difference between these romanization systems.

### Notice for Romaji Conversion
Since it's impossible to fully automatically convert __furigana__ directly to __romaji__ because furigana lacks information on pronunciation (Refer to [なぜ フリガナでは ダメなのか？](https://green.adam.ne.jp/roomazi/onamae.html#naze)). 

kuroshiro will not handle chōon when processing directly furigana (kana) -> romaji conversion with every romanization system (Except that Chōonpu will be handled) 

*For example, you'll get "kousi", "koushi", "koushi" respectively when converts kana "こうし" to romaji 
using `nippon`, `passport`, `hepburn` romanization system.*

The kanji -> romaji conversion with/without furigana mode is __unaffected__ by this logic.
