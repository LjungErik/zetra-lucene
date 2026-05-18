package filter_test

import (
	"testing"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/analysis/filter"
	"github.com/stretchr/testify/assert"
)

func Test_LowercaseFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		token    []analysis.Token
		expected []analysis.Token
	}{
		{
			name:     "lowercase already",
			token:    []analysis.Token{{Text: "testing", Position: 1}},
			expected: []analysis.Token{{Text: "testing", Position: 1}},
		},
		{
			name:     "capitilized",
			token:    []analysis.Token{{Text: "Testing", Position: 2}},
			expected: []analysis.Token{{Text: "testing", Position: 2}},
		},
		{
			name:     "all capitilized",
			token:    []analysis.Token{{Text: "FIRING", Position: 3}},
			expected: []analysis.Token{{Text: "firing", Position: 3}},
		},
		{
			name:     "mixed with special characters",
			token:    []analysis.Token{{Text: "Filter, this if YOU caN", Position: 4}},
			expected: []analysis.Token{{Text: "filter, this if you can", Position: 4}},
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

			f := filter.NewLowercaseFilter()
			actual := f.Apply(test.token)

			assert.Equal(t, test.expected, actual)
		})
	}
}
