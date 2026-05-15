package filter

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
)

type TrimFilter struct{}

var _ Filter = (*TrimFilter)(nil)

func NewTrimFilter() *TrimFilter {
	return &TrimFilter{}
}

func (f *TrimFilter) Apply(tokens []tokenizer.Token) []tokenizer.Token {
	n := 0
	for _, token := range tokens {
		trimmed := strings.TrimSpace(token.Text)
		if trimmed != "" {
			tokens[n] = token
			tokens[n].Text = trimmed
			n++
		}
	}
	return tokens[:n]
}
