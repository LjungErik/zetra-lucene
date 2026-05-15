package filter

import "github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"

type CompositeFilter struct {
	Filters []Filter
}

var _ Filter = (*CompositeFilter)(nil)

func (f *CompositeFilter) Apply(tokens []tokenizer.Token) []tokenizer.Token {
	for _, filter := range f.Filters {
		tokens = filter.Apply(tokens)
	}

	return tokens
}
