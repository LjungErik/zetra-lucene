package context

import (
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/search/score"
)

type SearchStatistics struct {
	AverageDataLength float64
	DocumentCount     int
}

type LeafReaderContext interface {
	GetTermCounts(fieldName string, term string) []index.TermCount
	GetStatistic(fieldName string) SearchStatistics
	GetDocLength(fieldName string, docID int) int
	GetSegmentID() int
}

type IndexReaderContext interface {
	GetSimilarity() score.Similarity
	GetLeaves() []LeafReaderContext
}
