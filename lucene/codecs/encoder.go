package codecs

type PostingsEncoder interface {
	StartTerm()
	StartDoc(docID, freq int) error
	AddPosition(pos int, p []byte) error
	FinishDoc()
	FinishTerm() error
}
