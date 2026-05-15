package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_LengthFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		min      int
		max      int
		token    []tokenizer.Token
		expected []tokenizer.Token
	}{
		{
			name:     "removes empty strings",
			min:      1,
			max:      256,
			token:    []tokenizer.Token{{Text: "", Position: 1}, {Text: "hello", Position: 5}},
			expected: []tokenizer.Token{{Text: "hello", Position: 5}},
		},
		{
			name:     "removes tokens below min length",
			min:      3,
			max:      256,
			token:    []tokenizer.Token{{Text: "ab", Position: 1}, {Text: "abc", Position: 4}, {Text: "a", Position: 8}},
			expected: []tokenizer.Token{{Text: "abc", Position: 4}},
		},
		{
			name:     "removes tokens above max length",
			min:      1,
			max:      5,
			token:    []tokenizer.Token{{Text: "hello", Position: 1}, {Text: "wonderful", Position: 7}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "removes tokens outside min and max range",
			min:      3,
			max:      5,
			token:    []tokenizer.Token{{Text: "ab", Position: 1}, {Text: "abc", Position: 4}, {Text: "abcdef", Position: 8}},
			expected: []tokenizer.Token{{Text: "abc", Position: 4}},
		},
		{
			name:     "keeps tokens at exact min and max boundaries",
			min:      3,
			max:      5,
			token:    []tokenizer.Token{{Text: "abc", Position: 1}, {Text: "abcde", Position: 5}},
			expected: []tokenizer.Token{{Text: "abc", Position: 1}, {Text: "abcde", Position: 5}},
		},
		{
			name:     "all tokens removed",
			min:      5,
			max:      10,
			token:    []tokenizer.Token{{Text: "ab", Position: 1}, {Text: "cd", Position: 4}},
			expected: []tokenizer.Token{},
		},
		{
			name:     "all tokens kept",
			min:      1,
			max:      256,
			token:    []tokenizer.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
			expected: []tokenizer.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
		},
		{
			name:     "empty input tokens",
			min:      1,
			max:      256,
			token:    []tokenizer.Token{},
			expected: []tokenizer.Token{},
		},
		{
			name:     "preserves positions of remaining tokens",
			min:      1,
			max:      256,
			token:    []tokenizer.Token{{Text: "", Position: 1}, {Text: "keep", Position: 5}, {Text: "", Position: 10}, {Text: "this", Position: 15}},
			expected: []tokenizer.Token{{Text: "keep", Position: 5}, {Text: "this", Position: 15}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f := filter.NewLengthFilter(test.min, test.max)
			actual := f.Apply(test.token)
			assert.Equal(t, test.expected, actual)
		})
	}
}
