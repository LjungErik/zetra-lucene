package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_LowercaseFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    tokenizer.Token
		expected []tokenizer.Token
	}{
		{
			name:     "lowercase already",
			token:    tokenizer.Token{Text: "testing", Position: 1},
			expected: []tokenizer.Token{{Text: "testing", Position: 1}},
		},
		{
			name:     "capitilized",
			token:    tokenizer.Token{Text: "Testing", Position: 2},
			expected: []tokenizer.Token{{Text: "testing", Position: 2}},
		},
		{
			name:     "all capitilized",
			token:    tokenizer.Token{Text: "FIRING", Position: 3},
			expected: []tokenizer.Token{{Text: "firing", Position: 3}},
		},
		{
			name:     "mixed with special characters",
			token:    tokenizer.Token{Text: "Filter, this if YOU caN", Position: 4},
			expected: []tokenizer.Token{{Text: "filter, this if you can", Position: 4}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			f := &filter.LowercaseFilter{}
			actual := f.Apply(test.token)

			assert.Equal(t, test.expected, actual)
		})
	}
}
