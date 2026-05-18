package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/stretchr/testify/assert"
)

func Test_StopWordFilter_EN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    []analysis.Token
		expected []analysis.Token
	}{
		{
			name:     "no stop word",
			token:    []analysis.Token{{Text: "testing", Position: 1}},
			expected: []analysis.Token{{Text: "testing", Position: 1}},
		},
		{
			name:     "only stop word #1",
			token:    []analysis.Token{{Text: "just", Position: 2}},
			expected: []analysis.Token{},
		},
		{
			name:     "only stop word #2",
			token:    []analysis.Token{{Text: "i", Position: 2}, {Text: "am", Position: 3}},
			expected: []analysis.Token{},
		},
		{
			name:     "stop words with space",
			token:    []analysis.Token{{Text: "i am just a man", Position: 3}},
			expected: []analysis.Token{{Text: "i am just a man", Position: 3}},
		},
		{
			name:     "no stop word",
			token:    []analysis.Token{{Text: "filter", Position: 4}},
			expected: []analysis.Token{{Text: "filter", Position: 4}},
		},
		{
			name:     "empty input",
			token:    []analysis.Token{},
			expected: []analysis.Token{},
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
