package lucene104

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShortestPrefix(t *testing.T) {
	tests := []struct {
		name     string
		first    string
		last     string
		prev     string
		expected int
	}{
		{
			name:     "Empty Prev test",
			first:    "monday",
			last:     "monkey",
			prev:     "",
			expected: 4,
		},
		{
			name:     "Example with Matching Prev",
			first:    "monkey",
			last:     "mothing",
			prev:     "monday",
			expected: 4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := shortestPrefixLength(test.prev, test.first, test.last)
			assert.Equal(t, test.expected, actual)
		})
	}
}
