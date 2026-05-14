package filter

import "github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"

type StopWordFilter struct {
	stopWords map[string]bool
}

func (f *StopWordFilter) Apply(token tokenizer.Token) []tokenizer.Token {
	if f.stopWords[token.Text] {
		return nil
	}

	return []tokenizer.Token{token}
}
