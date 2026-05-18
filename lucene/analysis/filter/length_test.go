package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/stretchr/testify/assert"
)

func Test_LengthFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		min      int
		max      int
		token    []analysis.Token
		expected []analysis.Token
	}{
		{
			name:     "removes empty strings",
			min:      1,
			max:      256,
			token:    []analysis.Token{{Text: "", Position: 1}, {Text: "hello", Position: 5}},
			expected: []analysis.Token{{Text: "hello", Position: 5}},
		},
		{
			name:     "removes tokens below min length",
			min:      3,
			max:      256,
			token:    []analysis.Token{{Text: "ab", Position: 1}, {Text: "abc", Position: 4}, {Text: "a", Position: 8}},
			expected: []analysis.Token{{Text: "abc", Position: 4}},
		},
		{
			name:     "removes tokens above max length",
			min:      1,
			max:      5,
			token:    []analysis.Token{{Text: "hello", Position: 1}, {Text: "wonderful", Position: 7}},
			expected: []analysis.Token{{Text: "hello", Position: 1}},
		},
		{
			name:     "removes tokens outside min and max range",
			min:      3,
			max:      5,
			token:    []analysis.Token{{Text: "ab", Position: 1}, {Text: "abc", Position: 4}, {Text: "abcdef", Position: 8}},
			expected: []analysis.Token{{Text: "abc", Position: 4}},
		},
		{
			name:     "keeps tokens at exact min and max boundaries",
			min:      3,
			max:      5,
			token:    []analysis.Token{{Text: "abc", Position: 1}, {Text: "abcde", Position: 5}},
			expected: []analysis.Token{{Text: "abc", Position: 1}, {Text: "abcde", Position: 5}},
		},
		{
			name:     "all tokens removed",
			min:      5,
			max:      10,
			token:    []analysis.Token{{Text: "ab", Position: 1}, {Text: "cd", Position: 4}},
			expected: []analysis.Token{},
		},
		{
			name:     "all tokens kept",
			min:      1,
			max:      256,
			token:    []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
			expected: []analysis.Token{{Text: "hello", Position: 1}, {Text: "world", Position: 7}},
		},
		{
			name:     "empty input tokens",
			min:      1,
			max:      256,
			token:    []analysis.Token{},
			expected: []analysis.Token{},
		},
		{
			name:     "preserves positions of remaining tokens",
			min:      1,
			max:      256,
			token:    []analysis.Token{{Text: "", Position: 1}, {Text: "keep", Position: 5}, {Text: "", Position: 10}, {Text: "this", Position: 15}},
			expected: []analysis.Token{{Text: "keep", Position: 5}, {Text: "this", Position: 15}},
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
