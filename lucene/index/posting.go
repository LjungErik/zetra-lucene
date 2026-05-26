package index

type Posting struct {
	freq       int
	documentID int
}

func (p *Posting) Frequency() int {
	return p.freq
}

func (p *Posting) DocID() int {
	return p.documentID
}
