package query

import (
	"github.com/LjungErik/zetra-lucene/lucene/search/collector"
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
)

type TermQuery struct {
	field string
	term  string
}

var _ Query = (*TermQuery)(nil)

func NewTermQuery(field, term string) *TermQuery {
	return &TermQuery{
		field: field,
		term:  term,
	}
}

func (q *TermQuery) Execute(ctx context.IndexReaderContext, col collector.TopDocumentCollector) {
	similarity := ctx.GetSimilarity()

	for _, leaf := range ctx.GetLeaves() {
		stats := leaf.GetStatistic(q.field)
		postings := leaf.GetTermCounts(q.field, q.term)

		scorer := similarity.GetScorer(stats.DocumentCount, len(postings), stats.AverageDataLength)

		for _, post := range postings {
			dl := leaf.GetDocLength(q.field, post.DocumentID)
			score := scorer(dl, post.Count)

			col.Add(&document.TopDoc{
				Score:      score,
				DocumentId: post.DocumentID,
				SegmentId:  leaf.GetSegmentID(),
			})
		}
	}
}
