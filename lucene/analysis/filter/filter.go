package filter

import "github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"

type Filter interface {
	Apply(tokenizer.Token) []tokenizer.Token
}
