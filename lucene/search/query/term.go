package query

import (
	"github.com/LjungErik/zetra-lucene/lucene/score"
	"github.com/LjungErik/zetra-lucene/lucene/search"
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

func (q *TermQuery) Execute(context search.SearchContext) []search.TopDoc {
	stats := context.GetStatistic(q.field)
	postings := context.GetTermCounts(q.field, q.term)

	topDocs := make([]search.TopDoc, len(postings))

	scorer := q.scoring.GetScorer(stats.DocumentCount, len(postings), stats.AverageDataLength)

	for i, post := range postings {
		dl := context.GetDocLength(q.field, post.DocumentID)
		score := scorer(dl)

		topDocs[i] = search.TopDoc{
			Score:      score,
			DocumentId: post.DocumentID,
			SegmentId:  context.GetSegmentID(),
		}
	}

	return topDocs
}
