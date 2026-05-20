package collector

import "github.com/LjungErik/zetra-lucene/lucene/search/document"

type TopDocumentCollector interface {
	Add(document.TopDoc)
}

type TopDocumentScoreCollector struct {
	limit int
}

var _ TopDocumentCollector = (*TopDocumentScoreCollector)(nil)

func NewTopDocumentScoreCollector(limit int) *TopDocumentScoreCollector {
	return &TopDocumentScoreCollector{
		limit: limit,
	}
}

func (c *TopDocumentScoreCollector) Add(doc document.TopDoc) {

}

func (c *TopDocumentScoreCollector) Collect() []document.TopDoc {
	return []document.TopDoc{}
}
