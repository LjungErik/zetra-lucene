package tokenizer

import "strings"

type DelimiterTokenizer struct {
	Delimiter string
}

var _ Tokenizer = (*DelimiterTokenizer)(nil)

func (t *DelimiterTokenizer) Tokenize(text string) []Token {
	parts := strings.Split(text, t.Delimiter)
	tokens := make([]Token, 0, len(parts))
	pos := 0
	for _, p := range parts {
		part := strings.TrimSpace(p)
		if part != "" {
			tokens = append(tokens, Token{Text: part, Position: pos})
			pos++
		}
	}

	return tokens
}
