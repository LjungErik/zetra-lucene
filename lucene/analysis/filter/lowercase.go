package filter

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
)

type LowercaseFilter struct{}

var _ Filter = (*LowercaseFilter)(nil)

func NewLowercaseFilter() *LowercaseFilter {
	return &LowercaseFilter{}
}

func (f *LowercaseFilter) Apply(tokens []tokenizer.Token) []tokenizer.Token {
	for i, token := range tokens {
		tokens[i].Text = strings.ToLower(token.Text)
	}

	return tokens
}
