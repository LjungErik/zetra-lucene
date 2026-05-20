package score

type Similarity interface {
	GetScorer(count int, postings int, avgDataLength float64) func(int, int) float64
}
