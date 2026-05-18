package search

import (
	"math"
	"sort"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/score"
	"github.com/LjungErik/zetra-lucene/lucene/storage"
)

type Document struct {
	DocumentID string
	Score      float64
}

type Query struct {
	Query         []analysis.Token
	Limit         int
	ScoreFunction *score.BM25Scoring
}

type Searcher interface {
	Search(query Query) []Document
}

type searcher struct {
	storage storage.Storage
}

var _ Searcher = (*searcher)(nil)

func NewSearcher(s storage.Storage) *searcher {
	return &searcher{
		storage: s,
	}
}

func (s *searcher) Search(query Query) []Document {
	scores := make(map[string]float64)
	n := float64(s.storage.TotalDocuments())
	avgDocLength := s.storage.GetAvgDocumentLength()

	for _, term := range query.Query {
		postings, ok := s.storage.GetPostings(term.Text)
		if !ok {
			continue
		}

		df := float64(len(postings))
		idf := math.Log((n-df+0.5)/(df+0.5) + 1.0)

		for _, posting := range postings {
			docLength := s.storage.GetDocumentLength(posting.DocumentID)
			scores[posting.DocumentID] += query.ScoreFunction.Score(posting.Count, docLength, avgDocLength, idf)
		}
	}

	foundDocs := make([]Document, 0, len(scores))
	for id, score := range scores {
		foundDocs = append(foundDocs, Document{
			Score:      score,
			DocumentID: id,
		})
	}

	sort.SliceStable(foundDocs, func(i, j int) bool {
		return foundDocs[i].Score > foundDocs[j].Score
	})

	limit := query.Limit
	if limit > len(foundDocs) {
		limit = len(foundDocs)
	}

	return foundDocs[:limit]
}
