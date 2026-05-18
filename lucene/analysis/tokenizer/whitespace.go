package tokenizer

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
)

type WhitespaceTokenizer struct{}

var _ Tokenizer = (*WhitespaceTokenizer)(nil)

func (t *WhitespaceTokenizer) Tokenize(text string) []analysis.Token {
	words := strings.Fields(text)

	tokens := make([]analysis.Token, 0, len(words))
	for i, word := range words {
		tokens = append(tokens, analysis.Token{Text: word, Position: i})
	}

	return tokens
}
