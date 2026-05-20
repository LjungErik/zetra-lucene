package document

type TopDoc struct {
	Score      float64
	DocumentId int
	SegmentId  int
}

func (a *TopDoc) Compare(b *TopDoc) int {
	if a.Score > b.Score {
		return -1
	} else if a.Score == b.Score {
		return 0
	}

	return 1
}
