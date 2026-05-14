package score

type BM25Scoring struct {
	K1 float64
	B  float64
}

func (s *BM25Scoring) Score(count, data_length int, avg_data_length, idf float64) float64 {
	dl := float64(data_length)
	c := float64(count)

	numerator := c * (s.K1 + 1.0)
	denominator := c + s.K1*(1.0-s.B+s.B*dl/avg_data_length)
	return idf * (numerator / denominator)
}
