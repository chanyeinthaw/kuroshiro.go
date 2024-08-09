package checker

import "testing"

func TestCheckers(t *testing.T) {
	// Hiragana
	if !IsHiragana('あ') {
		t.Errorf("IsHiragana('あ') should be true")
	}
	if IsHiragana('ア') {
		t.Errorf("IsHiragana('ア') should be false")
	}

	// Katakana
	if !IsKatakana('ア') {
		t.Errorf("IsKatakana('ア') should be true")
	}
	if IsKatakana('あ') {
		t.Errorf("IsKatakana('あ') should be false")
	}

	// Kana
	if !IsKana('あ') {
		t.Errorf("IsKana('あ') should be true")
	}
	if !IsKana('ア') {
		t.Errorf("IsKana('ア') should be true")
	}
	if IsKana('A') {
		t.Errorf("IsKana('A') should be false")
	}

	// Kanji
	if !IsKanji('漢') {
		t.Errorf("IsKanji('漢') should be true")
	}
	if IsKanji('あ') {
		t.Errorf("IsKanji('あ') should be false")
	}

	// Japanese
	if !IsJapanese('あ') {
		t.Errorf("IsJapanese('あ') should be true")
	}
	if !IsJapanese('漢') {
		t.Errorf("IsJapanese('漢') should be true")
	}
	if IsJapanese('A') {
		t.Errorf("IsJapanese('A') should be false")
	}
}

func TestHasKana(t *testing.T) {
	if !HasKana("あいうえお", HIRAGANA) {
		t.Errorf("HasKana('あいうえお', HIRAGANA) should be true")
	}
	if HasKana("あいうえお", KATAKANA) {
		t.Errorf("HasKana('あいうえお', KATAKANA) should be false")
	}
	if !HasKana("あいうえお", KANA) {
		t.Errorf("HasKana('あいうえお', KANA) should be true")
	}
	if HasKana("あいうえお", KANJI) {
		t.Errorf("HasKana('あいうえお', KANJI) should be false")
	}
	if !HasKana("あいうえお", JAPANESE) {
		t.Errorf("HasKana('あいうえお', JAPANESE) should be true")
	}
	if HasKana("abcde", HIRAGANA) {
		t.Errorf("HasKana('abcde', HIRAGANA) should be false")
	}
}

func TestGetStrType(t *testing.T) {
	if GetStrType("あいう漢字") != STR_KENJI_KANA {
		t.Errorf("GetStrType('あいう漢字') should be STR_KENJI_KANA")
	}
	if GetStrType("あいう") != STR_PURE_KANA {
		t.Errorf("GetStrType('あいう') should be STR_PURE_KANA")
	}
	if GetStrType("漢字") != STR_PURE_KENJI {
		t.Errorf("GetStrType('漢字') should be STR_PURE_KENJI")
	}
	if GetStrType("abcde") != STR_OTHER {
		t.Errorf("GetStrType('abcde') should be STR_OTHER")
	}
}
