package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/stretchr/testify/assert"
)

func Test_TrimFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    []analysis.Token
		expected []analysis.Token
	}{
		{
			name:     "trims leading whitespace",
			token:    []analysis.Token{{Text: "  hello", Position: 1}},
			expected: []analysis.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "trims trailing whitespace",
			token:    []analysis.Token{{Text: "hello  ", Position: 1}},
			expected: []analysis.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "trims leading and trailing whitespace",
			token:    []analysis.Token{{Text: "  hello  ", Position: 1}},
			expected: []analysis.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "does not trim whitespace within token",
			token:    []analysis.Token{{Text: "hello world", Position: 1}},
			expected: []analysis.Token{{Text: "hello world", Position: 1}},
		},
		{
			name:     "trims tabs and newlines",
			token:    []analysis.Token{{Text: "\thello\n", Position: 1}},
			expected: []analysis.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "space only token becomes empty string",
			token:    []analysis.Token{{Text: "   ", Position: 1}},
			expected: []analysis.Token{},
		},
		{
			name:     "whitespace only token becomes empty string",
			token:    []analysis.Token{{Text: "\n\n  \r\n \n  \t", Position: 1}},
			expected: []analysis.Token{},
		},
		{
			name:     "already trimmed tokens unchanged",
			token:    []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
			expected: []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
		},
		{
			name:     "trims across multiple tokens",
			token:    []analysis.Token{{Text: " hello ", Position: 1}, {Text: "\tworld\n", Position: 9}, {Text: "  foo  ", Position: 16}},
			expected: []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 9}, {Text: "foo", Position: 16}},
		},
		{
			name:     "empty input tokens",
			token:    []analysis.Token{},
			expected: []analysis.Token{},
		},
		{
			name:     "empty string token stays empty",
			token:    []analysis.Token{{Text: "", Position: 1}},
			expected: []analysis.Token{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f := filter.NewTrimFilter()
			actual := f.Apply(test.token)
			assert.Equal(t, test.expected, actual)
		})
	}
}
