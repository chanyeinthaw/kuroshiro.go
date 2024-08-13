package main

import (
	"fmt"

	"github.com/chanyeinthaw/kuroshiro.go"
	"github.com/chanyeinthaw/kuroshiro.go/analyzer"
)

const EXAMPLE_TEXT = "感じ取れたら手を繋ごう、重なるのは人生のライン and レミリア最高！"

func main() {
	analyzer, err := analyzer.NewMecab()
	defer analyzer.Destroy()
	if err != nil {
		panic(err)
	}

	input := EXAMPLE_TEXT

	ks := kuroshiro.New(analyzer)
	opts := kuroshiro.NewOptions().ConvertTo(kuroshiro.HIRAGANA).SetMode(kuroshiro.SPACED)
	result, err := ks.ConvertMultiline(input, opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(input)
	fmt.Println(result)
}
