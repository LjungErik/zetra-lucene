package filter_test

import (
	"regexp"
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/stretchr/testify/assert"
)

func Test_PatternReplaceFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pattern  *regexp.Regexp
		replace  string
		token    []analysis.Token
		expected []analysis.Token
	}{
		{
			name:     "remove parentheses",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  "",
			token:    []analysis.Token{{Text: "(hello)", Position: 1}, {Text: "world", Position: 8}},
			expected: []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 8}},
		},
		{
			name:     "token is only parentheses produces empty string",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  "",
			token:    []analysis.Token{{Text: "(())", Position: 1}, {Text: "keep", Position: 6}},
			expected: []analysis.Token{{Text: "", Position: 1}, {Text: "keep", Position: 6}},
		},
		{
			name:     "no match leaves tokens unchanged",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  "",
			token:    []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
			expected: []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
		},
		{
			name:     "replace with space instead of empty",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  " ",
			token:    []analysis.Token{{Text: "test(hello)", Position: 1}},
			expected: []analysis.Token{{Text: "test hello ", Position: 1}},
		},
		{
			name:     "remove digits from tokens",
			pattern:  regexp.MustCompile(`\d+`),
			replace:  "",
			token:    []analysis.Token{{Text: "abc123", Position: 1}, {Text: "456", Position: 8}},
			expected: []analysis.Token{{Text: "abc", Position: 1}, {Text: "", Position: 8}},
		},
		{
			name:     "replace hyphens with underscore",
			pattern:  regexp.MustCompile(`-`),
			replace:  "_",
			token:    []analysis.Token{{Text: "foo-bar-baz", Position: 1}},
			expected: []analysis.Token{{Text: "foo_bar_baz", Position: 1}},
		},
		{
			name:     "remove multiple special characters",
			pattern:  regexp.MustCompile(`[!@#$%^&*]`),
			replace:  "",
			token:    []analysis.Token{{Text: "h!e@l#l$o", Position: 1}, {Text: "w%o^r&l*d", Position: 12}},
			expected: []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 12}},
		},
		{
			name:     "empty input tokens",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  "",
			token:    []analysis.Token{},
			expected: []analysis.Token{},
		},
		{
			name:     "single token already empty",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  "",
			token:    []analysis.Token{{Text: "", Position: 1}},
			expected: []analysis.Token{{Text: "", Position: 1}},
		},
		{
			name:     "parentheses mixed with text across multiple tokens",
			pattern:  regexp.MustCompile(`[\(\)]`),
			replace:  "",
			token:    []analysis.Token{{Text: "func(a,b)", Position: 1}, {Text: "(test)", Position: 11}, {Text: "clean", Position: 18}},
			expected: []analysis.Token{{Text: "funca,b", Position: 1}, {Text: "test", Position: 11}, {Text: "clean", Position: 18}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f := filter.NewPatternReplaceFilter(test.pattern, test.replace)
			actual := f.Apply(test.token)
			assert.Equal(t, test.expected, actual)
		})
	}
}
