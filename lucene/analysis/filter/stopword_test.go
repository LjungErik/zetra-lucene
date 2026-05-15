package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_StopWordFilter_EN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    []tokenizer.Token
		expected []tokenizer.Token
	}{
		{
			name:     "no stop word",
			token:    []tokenizer.Token{{Text: "testing", Position: 1}},
			expected: []tokenizer.Token{{Text: "testing", Position: 1}},
		},
		{
			name:     "only stop word #1",
			token:    []tokenizer.Token{{Text: "just", Position: 2}},
			expected: []tokenizer.Token{},
		},
		{
			name:     "only stop word #2",
			token:    []tokenizer.Token{{Text: "i", Position: 2}, {Text: "am", Position: 3}},
			expected: []tokenizer.Token{},
		},
		{
			name:     "stop words with space",
			token:    []tokenizer.Token{{Text: "i am just a man", Position: 3}},
			expected: []tokenizer.Token{{Text: "i am just a man", Position: 3}},
		},
		{
			name:     "no stop word",
			token:    []tokenizer.Token{{Text: "filter", Position: 4}},
			expected: []tokenizer.Token{{Text: "filter", Position: 4}},
		},
		{
			name:     "empty input",
			token:    []tokenizer.Token{},
			expected: []tokenizer.Token{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f := filter.NewStopWordFilter(filter.EnglishStopWords)
			actual := f.Apply(test.token)

			assert.Equal(t, test.expected, actual)
		})
	}
}
