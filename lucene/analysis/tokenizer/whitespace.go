package tokenizer

import "strings"

type WhitespaceTokenizer struct{}

var _ Tokenizer = (*WhitespaceTokenizer)(nil)

func (t *WhitespaceTokenizer) Tokenize(text string) []Token {
	words := strings.Fields(text)

	tokens := make([]Token, 0, len(words))
	for i, word := range words {
		tokens = append(tokens, Token{Text: word, Position: i})
	}

	return tokens
}
