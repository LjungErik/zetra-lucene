package score

import "math"

type BM25Scoring struct {
	K1 float64
	B  float64
}

func (s *BM25Scoring) GetScorer(count int, postings int, avgDataLength float64) func(dataLength int) float64 {
	c := float64(count)
	df := float64(postings)
	idf := math.Log((c-df+0.5)/(df+0.5) + 0.1)

	return func(dataLength int) float64 {
		dl := float64(dataLength)
		numerator := c * (s.K1 + 1.0)
		denominator := c + s.K1*(1.0-s.B+s.B*dl/avgDataLength)

		return idf * (numerator / denominator)
	}
}
