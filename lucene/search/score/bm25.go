package score

import "math"

type BM25Similarity struct {
	K1 float64
	B  float64
}

var _ Similarity = (*BM25Similarity)(nil)

func NewBM25Similarity(k1, b float64) *BM25Similarity {
	return &BM25Similarity{
		K1: k1,
		B:  b,
	}
}

func (s *BM25Similarity) GetScorer(count int, postings int, avgDataLength float64) func(int, int) float64 {
	c := float64(count)
	df := float64(postings)
	idf := math.Log(1.0 + (c-df+0.5)/(df+0.5))

	return func(dataLength int, freq int) float64 {
		f := float64(freq)
		dl := float64(dataLength)
		numerator := f * (s.K1 + 1.0)
		denominator := f + s.K1*(1.0-s.B+s.B*dl/avgDataLength)

		return idf * (numerator / denominator)
	}
}
