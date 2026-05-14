package lucene

import "github.com/LjungErik/zetra-lucene/lucene/score"

type LuceneQuery struct {
	Query        string
	Total        int
	ScoreFuncion *score.BM25Scoring
}
