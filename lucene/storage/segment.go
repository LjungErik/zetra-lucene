package storage

import (
	"fmt"
	"math"
	"path/filepath"
	"sort"

	"github.com/LjungErik/zetra-lucene/lucene/score"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
)

type Segments struct {
	NextSegment  int64             `json:"next_segment"`
	SegmentCount int               `json:"segment_count"`
	Segments     []SegmentMetadata `json:"segments"`
}

type SegmentMetadata struct {
	SegmentName string `json:"segment_name"`
	SegmentID   int64  `json:"segment_id"`
}

type SegmentDocument struct {
	DocumentID string
	Score      float64
}

type SegmentQueryResults struct {
	Documents []SegmentDocument
}

type SegmentTermFreq struct {
	DocumentID string `json:"document_id"`
	Count      int    `json:"count"`
}

type SegmentDocumentsMetadata struct {
	DocsLength    map[string]int `json:"docs_length"`
	AvgDocsLength float64        `json:"avg_docs_length"`
}

type SegmentReader struct {
	metadata      SegmentMetadata
	index         map[string][]SegmentTermFreq
	docs          map[string]string
	docs_metadata SegmentDocumentsMetadata
	scoring       *score.BM25Scoring
}

func OpenSegmentReader(metadata SegmentMetadata, directory string) (*SegmentReader, error) {
	var (
		index_filename = filepath.Join(directory, fmt.Sprintf(metadata.SegmentName, ".idx"))
		docs_filename  = filepath.Join(directory, fmt.Sprintf(metadata.SegmentName, ".docs"))
		docs_metadata  = filepath.Join(directory, fmt.Sprintf(metadata.SegmentName, ".dmeta"))
	)

	r := &SegmentReader{
		metadata: metadata,
		scoring: &score.BM25Scoring{
			K1: 1.5,
			B:  0.75,
		},
	}

	if err := utils.ReadJsonFile(index_filename, &r.index); err != nil {
		return nil, err
	}

	if err := utils.ReadJsonFile(docs_filename, &r.docs); err != nil {
		return nil, err
	}

	if err := utils.ReadJsonFile(docs_metadata, &r.docs_metadata); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *SegmentReader) Query(term string, limit int) *SegmentQueryResults {
	postings, ok := r.index[term]
	n := len(r.docs)
	avgDocsLength := r.docs_metadata.AvgDocsLength
	if !ok {
		return nil
	}

	df := float64(len(postings))
	idf := math.Log((float64(n)-df+0.5)/(df+0.5) + 0.1)

	scores := make(map[string]float64, len(postings))

	for _, posting := range postings {
		docLength := r.docs_metadata.DocsLength[posting.DocumentID]
		scores[posting.DocumentID] += r.scoring.Score(n, docLength, avgDocsLength, idf)
	}

	docs := make([]SegmentDocument, 0, len(postings))
	for id, score := range scores {
		docs = append(docs, SegmentDocument{
			DocumentID: id,
			Score:      score,
		})
	}

	sort.SliceStable(docs, func(i, j int) bool {
		return docs[i].Score > docs[j].Score
	})

	if limit > len(docs) {
		return &SegmentQueryResults{
			Documents: docs,
		}
	}

	return &SegmentQueryResults{
		Documents: docs[:limit],
	}

}
