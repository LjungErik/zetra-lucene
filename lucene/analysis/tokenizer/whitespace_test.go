package tokenizer_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_Whitespace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		text     string
		expected []tokenizer.Token
	}{
		{
			name: "phase with only space",
			text: "I am just testing this phrase",
			expected: []tokenizer.Token{
				{Text: "I", Position: 0},
				{Text: "am", Position: 1},
				{Text: "just", Position: 2},
				{Text: "testing", Position: 3},
				{Text: "this", Position: 4},
				{Text: "phrase", Position: 5},
			},
		},
		{
			name: "phase with space, newline and tab",
			text: "I am\r\r\njust testing\nthis phrase\t and this as well.",
			expected: []tokenizer.Token{
				{Text: "I", Position: 0},
				{Text: "am", Position: 1},
				{Text: "just", Position: 2},
				{Text: "testing", Position: 3},
				{Text: "this", Position: 4},
				{Text: "phrase", Position: 5},
				{Text: "and", Position: 6},
				{Text: "this", Position: 7},
				{Text: "as", Position: 8},
				{Text: "well.", Position: 9},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			tok := tokenizer.WhitespaceTokenizer{}
			actual := tok.Tokenize(test.text)
			assert.Equal(t, test.expected, actual)
		})
	}
}
