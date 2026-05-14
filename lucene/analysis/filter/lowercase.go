package filter

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
)

type LowercaseFilter struct{}

var _ Filter = (*LowercaseFilter)(nil)

func (f *LowercaseFilter) Apply(token tokenizer.Token) []tokenizer.Token {
	return []tokenizer.Token{{Text: strings.ToLower(token.Text), Position: token.Position}}
}
