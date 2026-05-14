package tokenizer

import "strings"

type CompositeTokenizer struct {
	Tokenizers []Tokenizer
}

var _ Tokenizer = (*CompositeTokenizer)(nil)

func (t *CompositeTokenizer) Tokenize(text string) []Token {
	texts := []string{text}

	for _, tok := range t.Tokenizers {
		var next []string

		for _, text := range texts {
			for _, token := range tok.Tokenize(text) {
				next = append(next, token.Text)
			}
		}

		texts = next
	}

	tokens := make([]Token, 0, len(texts))
	pos := 0
	for _, text := range texts {
		trimmed := strings.TrimSpace(text)
		if trimmed != "" {
			tokens = append(tokens, Token{Text: trimmed, Position: pos})
			pos++
		}
	}

	return tokens
}
