package lucene_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/LjungErik/zetra-lucene/lucene"
)

type TestDocument struct {
	DocumentID string
	Data       string
}

func Test_LuceneIndex_Search(t *testing.T) {
	tests := []struct {
		name         string
		docs         []TestDocument
		query        string
		total        int
		k1           float64
		b            float64
		expectedDocs []TestDocument
	}{
		{
			name:  "test simple query",
			query: "How many bones are there in a fish?",
			total: 10,
			k1:    1.5,
			b:     0.75,
			docs: []TestDocument{
				{
					DocumentID: "doc-1",
					Data:       "magic document starts here",
				},
				{
					DocumentID: "doc-2",
					Data:       "A fish has on average 124 bones",
				},
				{
					DocumentID: "doc-3",
					Data:       "A human has 207 bones",
				},
			},
			expectedDocs: []TestDocument{
				{
					DocumentID: "doc-2",
					Data:       "A fish has on average 124 bones",
				},
				{
					DocumentID: "doc-3",
					Data:       "A human has 207 bones",
				},
			},
		},
		{
			name:  "test query with multiple matches",
			query: "Are fish good at flying?",
			total: 2,
			k1:    1.5,
			b:     0.75,
			docs: []TestDocument{
				{
					DocumentID: "doc-1",
					Data:       "magic document starts here\nand ends after we have found all the fish",
				},
				{
					DocumentID: "doc-2",
					Data:       "A fish has on average 124 bones",
				},
				{
					DocumentID: "doc-3",
					Data:       "A human has 207 bones",
				},
				{
					DocumentID: "doc-4",
					Data:       "Fish are great at flying but they don't really survive for long outside of water",
				},
				{
					DocumentID: "doc-5",
					Data:       "one fish, two fish, three fish, gold fish",
				},
			},
			expectedDocs: []TestDocument{
				{
					DocumentID: "doc-4",
					Data:       "Fish are great at flying but they don't really survive for long outside of water",
				},
				{
					DocumentID: "doc-5",
					Data:       "one fish, two fish, three fish, gold fish",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			index := lucene.NewIndex()
			for _, doc := range test.docs {
				err := index.Add(doc.DocumentID, doc.Data)

				assert.NoError(t, err)
			}

			docs := index.Search(test.query, test.k1, test.b, test.total)
			assert.Equal(t, len(test.expectedDocs), len(docs))

			for i, doc := range docs {
				assert.Equal(t, test.expectedDocs[i].DocumentID, doc.DocumentID)
				assert.Equal(t, test.expectedDocs[i].Data, doc.Document)
			}
		})
	}
}
