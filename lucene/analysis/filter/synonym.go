package filter

import "github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"

type SynonymFilter struct {
	synonyms map[string][]string
}

func (f *SynonymFilter) Apply(token tokenizer.Token) []tokenizer.Token {
	result := []tokenizer.Token{token}
	for _, syn := range f.synonyms[token.Text] {
		result = append(result, tokenizer.Token{Text: syn, Position: token.Position})
	}

	return result
}
