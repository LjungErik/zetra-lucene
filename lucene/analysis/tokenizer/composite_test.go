package tokenizer_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_Composite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		text      string
		tokenizer *tokenizer.CompositeTokenizer
		expected  []tokenizer.Token
	}{
		{
			name: "phase with space and dots",
			text: "I am just testing this phrase. More text follows after the dot.",
			tokenizer: &tokenizer.CompositeTokenizer{
				Tokenizers: []tokenizer.Tokenizer{
					&tokenizer.WhitespaceTokenizer{},
					&tokenizer.DelimiterTokenizer{Delimiter: "."},
				},
			},
			expected: []tokenizer.Token{
				{Text: "I", Position: 0},
				{Text: "am", Position: 1},
				{Text: "just", Position: 2},
				{Text: "testing", Position: 3},
				{Text: "this", Position: 4},
				{Text: "phrase", Position: 5},
				{Text: "More", Position: 6},
				{Text: "text", Position: 7},
				{Text: "follows", Position: 8},
				{Text: "after", Position: 9},
				{Text: "the", Position: 10},
				{Text: "dot", Position: 11},
			},
		},
		{
			name: "phase with whitespace and special delimiter characters",
			text: "I am just testing\nthis phrase, and this as well. Lets see if this magic-works for,me,you,and I.",
			tokenizer: &tokenizer.CompositeTokenizer{
				Tokenizers: []tokenizer.Tokenizer{
					&tokenizer.WhitespaceTokenizer{},
					&tokenizer.DelimiterTokenizer{Delimiter: "."},
					&tokenizer.DelimiterTokenizer{Delimiter: ","},
					&tokenizer.DelimiterTokenizer{Delimiter: "-"},
				},
			},
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
				{Text: "well", Position: 9},
				{Text: "Lets", Position: 10},
				{Text: "see", Position: 11},
				{Text: "if", Position: 12},
				{Text: "this", Position: 13},
				{Text: "magic", Position: 14},
				{Text: "works", Position: 15},
				{Text: "for", Position: 16},
				{Text: "me", Position: 17},
				{Text: "you", Position: 18},
				{Text: "and", Position: 19},
				{Text: "I", Position: 20},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := test.tokenizer.Tokenize(test.text)
			assert.Equal(t, test.expected, actual)
		})
	}
}
