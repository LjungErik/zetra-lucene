package lucene_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/analyzer"
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/document/field/textfield"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/index/writer"
	"github.com/LjungErik/zetra-lucene/lucene/search"
	searchdoc "github.com/LjungErik/zetra-lucene/lucene/search/document"
	"github.com/LjungErik/zetra-lucene/lucene/search/query"
	"github.com/LjungErik/zetra-lucene/lucene/search/reader"
)

type TestDocument struct {
	Field string
	Data  string
}

func Test_Indexing_and_Search(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		docs         []TestDocument
		query        query.Query
		total        int
		expectedDocs []searchdoc.TopDoc
	}{
		{
			name:  "test simple query",
			query: query.NewTermQuery("name", "fish"),
			total: 10,
			docs: []TestDocument{
				{
					Field: "name",
					Data:  "magic document starts here",
				},
				{
					Field: "name",
					Data:  "A fish has on average 124 bones",
				},
				{
					Field: "name",
					Data:  "A human has 207 bones",
				},
			},
			expectedDocs: []searchdoc.TopDoc{
				{DocumentId: 1, SegmentId: 0},
			},
		},
		{
			name:  "test query with multiple matches",
			query: query.NewTermQuery("name", "fish"),
			total: 2,
			docs: []TestDocument{
				{
					Field: "name",
					Data:  "magic document starts here\nand ends after we have found all the fish",
				},
				{
					Field: "name",
					Data:  "A fish has on average 124 bones",
				},
				{
					Field: "name",
					Data:  "A human has 207 bones",
				},
				{
					Field: "name",
					Data:  "Fish are great at flying but they don't really survive for long outside of water",
				},
				{
					Field: "name",
					Data:  "one fish, two fish, three fish, gold fish",
				},
			},
			expectedDocs: []searchdoc.TopDoc{
				{DocumentId: 4, SegmentId: 0},
				{DocumentId: 1, SegmentId: 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//t.Parallel()
			tmpDir := t.TempDir()
			defer os.RemoveAll(tmpDir)

			dir := directory.OpenFSDirectory(tmpDir)
			fieldAnalyzer := analyzer.NewPerFieldAnalyzer()

			indexWriter := writer.NewIndexWriter(dir, writer.IndexWriterConfig{
				Analyzer: fieldAnalyzer,
			})

			for _, testDoc := range test.docs {
				doc := document.NewDocument()
				doc.Add(textfield.New(testDoc.Field, testDoc.Data, true))
				indexWriter.AddDocument(doc)
			}

			err := indexWriter.Flush()
			require.NoError(t, err)

			dirReader, err := reader.OpenStandrardDirectoryReader(dir)
			require.NoError(t, err)

			searcher := search.NewIndexSearcher(dirReader)

			topDocs := searcher.Query(test.query, test.total)
			require.Equal(t, len(test.expectedDocs), len(topDocs))

			for i, topDoc := range topDocs {
				assert.Equal(t, test.expectedDocs[i].DocumentId, topDoc.DocumentId)
				assert.Equal(t, test.expectedDocs[i].SegmentId, topDoc.SegmentId)
			}
		})
	}
}
