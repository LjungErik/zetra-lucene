package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/tokenizer"
	"github.com/stretchr/testify/assert"
)

func Test_CompositeFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filter   *filter.CompositeFilter
		token    []tokenizer.Token
		expected []tokenizer.Token
	}{
		{
			name: "lowercase already",
			filter: &filter.CompositeFilter{
				Filters: []filter.Filter{filter.NewLowercaseFilter()},
			},
			token:    []tokenizer.Token{{Text: "TeStiNg", Position: 1}, {Text: "AgaiN", Position: 12}},
			expected: []tokenizer.Token{{Text: "testing", Position: 1}, {Text: "again", Position: 12}},
		},
		{
			name: "capitilized with stop words filter",
			filter: &filter.CompositeFilter{
				Filters: []filter.Filter{filter.NewLowercaseFilter(), filter.NewStopWordFilter(filter.EnglishStopWords)},
			},
			token:    []tokenizer.Token{{Text: "Testing", Position: 2}, {Text: "IN", Position: 4}, {Text: "A", Position: 3}},
			expected: []tokenizer.Token{{Text: "testing", Position: 2}},
		},
		{
			name: "all capitilized without stop words filter",
			filter: &filter.CompositeFilter{
				Filters: []filter.Filter{filter.NewLowercaseFilter()},
			},
			token:    []tokenizer.Token{{Text: "FIRING", Position: 3}, {Text: "IN", Position: 4}, {Text: "A", Position: 3}},
			expected: []tokenizer.Token{{Text: "firing", Position: 3}, {Text: "in", Position: 4}, {Text: "a", Position: 3}},
		},
		{
			name: "Only stop words filter",
			filter: &filter.CompositeFilter{
				Filters: []filter.Filter{filter.NewStopWordFilter(filter.EnglishStopWords)},
			},
			token:    []tokenizer.Token{{Text: "FIRING in a HOUSE", Position: 5}, {Text: "in", Position: 4}, {Text: "a", Position: 3}},
			expected: []tokenizer.Token{{Text: "FIRING in a HOUSE", Position: 5}},
		},
		{
			name: "empty input",
			filter: &filter.CompositeFilter{
				Filters: []filter.Filter{filter.NewLowercaseFilter(), filter.NewStopWordFilter(filter.EnglishStopWords)},
			},
			token:    []tokenizer.Token{},
			expected: []tokenizer.Token{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := test.filter.Apply(test.token)

			assert.Equal(t, test.expected, actual)
		})
	}
}
