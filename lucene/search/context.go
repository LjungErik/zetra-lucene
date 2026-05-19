package search

import "github.com/LjungErik/zetra-lucene/lucene/index"

type SearchStatistics struct {
	AverageDataLength float64
	DocumentCount     int
}

type SearchContext interface {
	GetTermCounts(fieldName string, term string) []index.TermCount
	GetStatistic(fieldName string) SearchStatistics
	GetDocLength(fieldName string, docID int) int
	GetSegmentID() int
}
