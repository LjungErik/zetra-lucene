package filter

import (
	"regexp"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
)

type PatternReplaceFilter struct {
	pattern     *regexp.Regexp
	replacement string
}

var _ Filter = (*PatternReplaceFilter)(nil)

func NewPatternReplaceFilter(pattern *regexp.Regexp, replacement string) *PatternReplaceFilter {
	return &PatternReplaceFilter{
		pattern:     pattern,
		replacement: replacement,
	}
}

func (f *PatternReplaceFilter) Apply(tokens []analysis.Token) []analysis.Token {
	for i := range tokens {
		tokens[i].Text = f.pattern.ReplaceAllString(tokens[i].Text, f.replacement)
	}

	return tokens
}
