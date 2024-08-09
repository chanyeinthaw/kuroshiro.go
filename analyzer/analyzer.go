package analyzer

type RawResult [][]string

type Token struct {
	SurfaceForm    string
	Pos            string
	PosDetail1     string
	PosDetail2     string
	PosDetail3     string
	ConjugatedType string
	ConjugatedForm string
	BasicForm      string
	Reading        string
	Pronunciation  string
}

const (
	NUL = "\x00"
)

type Analyzer interface {
	Destroy()
	ParseRaw(text string) (RawResult, error)
	Parse(text string) ([]Token, error)
}
