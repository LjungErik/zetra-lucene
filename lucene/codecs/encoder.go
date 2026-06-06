package codecs

type PostingsEncoder interface {
	StartTerm()
	StartDoc(docID, freq int)
	AddPosition(pos int, p []byte)
	FinishDoc()
	FinishTerm()
}
