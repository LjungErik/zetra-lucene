package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_TrimFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    []tokenizer.Token
		expected []tokenizer.Token
	}{
		{
			name:     "trims leading whitespace",
			token:    []tokenizer.Token{{Text: "  hello", Position: 1}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "trims trailing whitespace",
			token:    []tokenizer.Token{{Text: "hello  ", Position: 1}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "trims leading and trailing whitespace",
			token:    []tokenizer.Token{{Text: "  hello  ", Position: 1}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "does not trim whitespace within token",
			token:    []tokenizer.Token{{Text: "hello world", Position: 1}},
			expected: []tokenizer.Token{{Text: "hello world", Position: 1}},
		},
		{
			name:     "trims tabs and newlines",
			token:    []tokenizer.Token{{Text: "\thello\n", Position: 1}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "space only token becomes empty string",
			token:    []tokenizer.Token{{Text: "   ", Position: 1}},
			expected: []tokenizer.Token{},
		},
		{
			name:     "whitespace only token becomes empty string",
			token:    []tokenizer.Token{{Text: "\n\n  \r\n \n  \t", Position: 1}},
			expected: []tokenizer.Token{},
		},
		{
			name:     "already trimmed tokens unchanged",
			token:    []tokenizer.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
		},
		{
			name:     "trims across multiple tokens",
			token:    []tokenizer.Token{{Text: " hello ", Position: 1}, {Text: "\tworld\n", Position: 9}, {Text: "  foo  ", Position: 16}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 9}, {Text: "foo", Position: 16}},
		},
		{
			name:     "empty input tokens",
			token:    []tokenizer.Token{},
			expected: []tokenizer.Token{},
		},
		{
			name:     "empty string token stays empty",
			token:    []tokenizer.Token{{Text: "", Position: 1}},
			expected: []tokenizer.Token{},
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
