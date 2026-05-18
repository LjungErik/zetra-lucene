package filter

import "github.com/LjungErik/zetra-lucene/lucene/analysis"

type LengthFilter struct {
	min int
	max int
}

var _ Filter = (*LengthFilter)(nil)

func NewLengthFilter(min, max int) *LengthFilter {
	return &LengthFilter{
		min: min,
		max: max,
	}
}

func (f *LengthFilter) Apply(tokens []analysis.Token) []analysis.Token {
	n := 0
	for _, token := range tokens {
		if len(token.Text) >= f.min && len(token.Text) <= f.max {
			tokens[n] = token
			n++
		}
	}

	return tokens[:n]
}
