package tokenizer

import "github.com/LjungErik/zetra-lucene/lucene/analysis"

type Tokenizer interface {
	Tokenize(string) []analysis.Token
}
