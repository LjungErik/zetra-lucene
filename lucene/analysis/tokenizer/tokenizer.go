package tokenizer

type Token struct {
	Text     string
	Position int
}

type Tokenizer interface {
	Tokenize(string) []Token
}
