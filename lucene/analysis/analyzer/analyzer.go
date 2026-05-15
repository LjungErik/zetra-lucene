package analyzer

import (
	"regexp"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
)

type Analyzer interface {
	Analyze(data string) []tokenizer.Token
}

type DefaultAnalyzer struct {
	tokenizer tokenizer.Tokenizer
	filter    filter.Filter
}

var _ Analyzer = (*DefaultAnalyzer)(nil)

func (a *DefaultAnalyzer) Analyze(data string) []tokenizer.Token {
	tokens := a.tokenizer.Tokenize(data)
	tokens = a.filter.Apply(tokens)

	return tokens
}

func NewEnglishLanguageAnalyzer() Analyzer {
	return &DefaultAnalyzer{
		tokenizer: &tokenizer.CompositeTokenizer{
			Tokenizers: []tokenizer.Tokenizer{
				&tokenizer.WhitespaceTokenizer{},
				&tokenizer.DelimiterTokenizer{Delimiter: "."},
				&tokenizer.DelimiterTokenizer{Delimiter: ","},
				&tokenizer.DelimiterTokenizer{Delimiter: "-"},
				&tokenizer.DelimiterTokenizer{Delimiter: "_"},
				&tokenizer.DelimiterTokenizer{Delimiter: ";"},
			},
		},
		filter: &filter.CompositeFilter{
			Filters: []filter.Filter{
				filter.NewLowercaseFilter(),
				filter.NewStopWordFilter(filter.EnglishStopWords),
				filter.NewPatternReplaceFilter(regexp.MustCompile(`[\(\)]`), ""),
				filter.NewPatternReplaceFilter(regexp.MustCompile(`[!@#$%^&*]`), ""),
				filter.NewTrimFilter(),
			},
		},
	}
}
