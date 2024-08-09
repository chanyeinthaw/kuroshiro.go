package analyzer

import "testing"

const EXAMPLE_TEXT = "すもももももも"

func TestInit(t *testing.T) {
	a, err := NewMecab()
	defer a.Destroy()

	if err != nil {
		t.Errorf("Failed to initialize mecab")
	}
}

func TestParseSentance(t *testing.T) {
	analyzer, _ := NewMecab()
	defer analyzer.Destroy()
	result, err := analyzer.Parse(EXAMPLE_TEXT)
	if err != nil {
		t.Errorf("Failed to parse text err %v", err)
	}
	if len(result) != 4 {
		t.Errorf("Failed to parse text. Invalid result array length")
	}
}

func TestParseBlank(t *testing.T) {
	analyzer, _ := NewMecab()
	defer analyzer.Destroy()
	result, err := analyzer.Parse("")
	if err != nil {
		t.Errorf("Failed to parse text err %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Failed to parse text. Invalid result array length")
	}
}
