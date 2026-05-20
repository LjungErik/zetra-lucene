package query

import (
	"github.com/LjungErik/zetra-lucene/lucene/score"
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
)

type TermQuery struct {
	field   string
	term    string
	scoring *score.BM25Scoring
}

var _ Query = (*TermQuery)(nil)

func NewTermQuery(field, term string) *TermQuery {
	return &TermQuery{
		field: field,
		term:  term,
		scoring: &score.BM25Scoring{
			K1: 1.5,
			B:  0.75,
		},
	}
}

func (q *TermQuery) Execute(ctx context.SearchContext) []document.TopDoc {
	stats := ctx.GetStatistic(q.field)
	postings := ctx.GetTermCounts(q.field, q.term)

	topDocs := make([]document.TopDoc, len(postings))

	scorer := q.scoring.GetScorer(stats.DocumentCount, len(postings), stats.AverageDataLength)

	for i, post := range postings {
		dl := ctx.GetDocLength(q.field, post.DocumentID)
		score := scorer(dl, post.Count)

		topDocs[i] = document.TopDoc{
			Score:      score,
			DocumentId: post.DocumentID,
			SegmentId:  ctx.GetSegmentID(),
		}
	}

	return topDocs
}
