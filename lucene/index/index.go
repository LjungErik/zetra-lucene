package index

import (
	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/storage"
)

type Document struct {
	DocumentID string
	Data       string
	Tokens     []analysis.Token
}

type Indexer interface {
	Index(document *Document) error
}

type indexer struct {
	storage storage.Storage
}

var _ Indexer = (*indexer)(nil)

func NewIndexer(s storage.Storage) *indexer {
	return &indexer{
		storage: s,
	}
}

func (i *indexer) Index(document *Document) error {
	// Current implementation we don't enable overwrite or update
	if i.storage.IsIndexed(document.DocumentID) {
		return ErrDocumentAlreadyIndexed
	}

	term_count := make(map[string]int, len(document.Tokens))

	for _, token := range document.Tokens {
		term_count[token.Text]++
	}

	terms := make([]storage.DocumentTermCount, 0, len(term_count))
	for term, count := range term_count {
		terms = append(terms, storage.DocumentTermCount{
			Term:  term,
			Count: count,
		})
	}

	i.storage.Insert(storage.Document{
		DocumentID: document.DocumentID,
		TermCounts: terms,
	})

	i.storage.InsertDocument(document.DocumentID, document.Data)

	return nil
}
