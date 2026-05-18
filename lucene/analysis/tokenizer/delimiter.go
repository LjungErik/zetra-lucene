package tokenizer

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
)

type DelimiterTokenizer struct {
	Delimiter string
}

var _ Tokenizer = (*DelimiterTokenizer)(nil)

func (t *DelimiterTokenizer) Tokenize(text string) []analysis.Token {
	parts := strings.Split(text, t.Delimiter)
	tokens := make([]analysis.Token, 0, len(parts))
	pos := 0
	for _, p := range parts {
		part := strings.TrimSpace(p)
		if part != "" {
			tokens = append(tokens, analysis.Token{Text: part, Position: pos})
			pos++
		}
	}

	return tokens
}
