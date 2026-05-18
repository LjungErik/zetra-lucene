package tokenizer_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_Delimiter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		text      string
		delimiter string
		expected  []analysis.Token
	}{
		{
			name:      "phase with single comma",
			text:      "I am just testing this phrase, and having a new phase after.",
			delimiter: ",",
			expected: []analysis.Token{
				{Text: "I am just testing this phrase", Position: 0},
				{Text: "and having a new phase after.", Position: 1},
			},
		},
		{
			name:      "phase with multiple comma",
			text:      "Mix the eggs, butter, sugar, and flour in a bowl.",
			delimiter: ",",
			expected: []analysis.Token{
				{Text: "Mix the eggs", Position: 0},
				{Text: "butter", Position: 1},
				{Text: "sugar", Position: 2},
				{Text: "and flour in a bowl.", Position: 3},
			},
		},
		{
			name:      "phase with dots",
			text:      "I wonder... will this work...\nor maybe not.",
			delimiter: ".",
			expected: []analysis.Token{
				{Text: "I wonder", Position: 0},
				{Text: "will this work", Position: 1},
				{Text: "or maybe not", Position: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tok := tokenizer.DelimiterTokenizer{Delimiter: test.delimiter}
			actual := tok.Tokenize(test.text)
			assert.Equal(t, test.expected, actual)
		})
	}
}
