package search

import (
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
	"github.com/LjungErik/zetra-lucene/lucene/search/query"
	"github.com/LjungErik/zetra-lucene/lucene/search/reader"
)

type IndexSearcher struct {
	reader reader.DirectoryReader
}

func NewIndexSearcher(r reader.DirectoryReader) *IndexSearcher {
	return &IndexSearcher{
		reader: r,
	}
}

func (s *IndexSearcher) Query(q query.Query, n int) []document.TopDoc {
	return s.reader.Query(q, n)
}
