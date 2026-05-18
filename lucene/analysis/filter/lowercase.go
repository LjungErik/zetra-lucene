package filter

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
)

type LowercaseFilter struct{}

var _ Filter = (*LowercaseFilter)(nil)

func NewLowercaseFilter() *LowercaseFilter {
	return &LowercaseFilter{}
}

func (f *LowercaseFilter) Apply(tokens []analysis.Token) []analysis.Token {
	for i, token := range tokens {
		tokens[i].Text = strings.ToLower(token.Text)
	}

	return tokens
}
