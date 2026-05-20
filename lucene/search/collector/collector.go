package collector

import (
	"github.com/LjungErik/zetra-lucene/lucene/internal/queue"
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
)

type TopDocumentCollector interface {
	Add(*document.TopDoc)
}

type TopDocumentScoreCollector struct {
	limit int
	queue *queue.PriorityQueue[*document.TopDoc]
}

var _ TopDocumentCollector = (*TopDocumentScoreCollector)(nil)

func NewTopDocumentScoreCollector(limit int) *TopDocumentScoreCollector {
	return &TopDocumentScoreCollector{
		limit: limit,
		queue: queue.NewPriorityQueue[*document.TopDoc](limit),
	}
}

func (c *TopDocumentScoreCollector) Add(doc *document.TopDoc) {
	c.queue.Push(doc)
}

func (c *TopDocumentScoreCollector) Collect() []*document.TopDoc {
	limit := c.limit
	if limit > c.queue.Len() {
		limit = c.queue.Len()
	}

	ret := make([]*document.TopDoc, limit)
	for i := 0; i < limit; i++ {
		ret[i], _ = c.queue.Pop()
	}

	return ret
}
