package analysis

import (
	"iter"
)

type Token struct {
	Text     string
	Position int
}

type TokenStream struct {
	tokens []Token
}

func NewTokenStream(tokens []Token) *TokenStream {
	return &TokenStream{
		tokens: tokens,
	}
}

func (s *TokenStream) Iter() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		for _, tok := range s.tokens {
			if !yield(tok) {
				return
			}
		}
	}
}
