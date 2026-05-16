package storage

type Document struct {
	DocumentID string
	TermCounts []DocumentTermCount
}

type DocumentTermCount struct {
	Term  string
	Count int
}

type Storage interface {
	IsIndexed(documentID string) bool
	Insert(document Document)
	InsertDocument(documentID string, data string)
	TotalDocuments() int
	GetPostings(term string) ([]TermFreqency, bool)
	GetDocumentLength(documentID string) int
	GetAvgDocumentLength() float64
	GetDocument(documentID string) string
}

type TermFreqency struct {
	DocumentID string
	Count      int
}

type storage struct {
	index          map[string][]TermFreqency
	docs           map[string]string
	doc_length     map[string]int
	avg_doc_length float64
}

var _ Storage = (*storage)(nil)

func NewStorage() *storage {
	return &storage{
		index:          make(map[string][]TermFreqency),
		docs:           make(map[string]string),
		doc_length:     make(map[string]int),
		avg_doc_length: 0.0,
	}
}

func (s *storage) IsIndexed(documentID string) bool {
	_, ok := s.doc_length[documentID]

	return ok
}

func (s *storage) Insert(document Document) {
	for _, term := range document.TermCounts {
		s.index[term.Term] = append(s.index[term.Term], TermFreqency{
			DocumentID: document.DocumentID,
			Count:      term.Count,
		})
	}
}

func (s *storage) InsertDocument(documentID string, data string) {
	s.docs[documentID] = data

	s.doc_length[documentID] = len(data)

	sum := 0
	for _, length := range s.doc_length {
		sum += length
	}

	s.avg_doc_length = float64(sum) / float64(len(s.doc_length))
}

func (s *storage) TotalDocuments() int {
	return len(s.docs)
}

func (s *storage) GetPostings(term string) ([]TermFreqency, bool) {
	freq, ok := s.index[term]
	return freq, ok
}

func (s *storage) GetDocumentLength(documentID string) int {
	return s.doc_length[documentID]
}

func (s *storage) GetAvgDocumentLength() float64 {
	return s.avg_doc_length
}

func (s *storage) GetDocument(documentID string) string {
	return s.docs[documentID]
}
