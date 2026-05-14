package filter

import (
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
)

type LowercaseFilter struct{}

func (f *LowercaseFilter) Apply(token tokenizer.Token) []tokenizer.Token {
	return []tokenizer.Token{{Text: strings.ToLower(token.Text), Position: token.Position}}
}
