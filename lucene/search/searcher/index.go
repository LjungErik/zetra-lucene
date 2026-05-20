package search

import (
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
	"github.com/LjungErik/zetra-lucene/lucene/search/query"
	"github.com/LjungErik/zetra-lucene/lucene/search/query/collector"
	"github.com/LjungErik/zetra-lucene/lucene/search/reader"
	"github.com/LjungErik/zetra-lucene/lucene/search/score"
)

type IndexSearcher struct {
	reader     reader.DirectoryReader
	similarity score.Similarity
}

type Option func(*IndexSearcher)

var _ context.IndexReaderContext = (*IndexSearcher)(nil)

func NewIndexSearcher(r reader.DirectoryReader, opts ...Option) *IndexSearcher {
	s := &IndexSearcher{
		reader: r,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *IndexSearcher) Query(q query.Query, n int) []document.TopDoc {
	col := collector.NewTopDocumentScoreCollector(n)

	q.Execute(s, col)

	return col.Collect()
}

func (s *IndexSearcher) GetLeaves() []context.LeafReaderContext {
	return s.reader.GetSegments()
}

func (s *IndexSearcher) GetSimilarity() score.Similarity {
	return s.similarity
}
