package lucene

type LuceneIndexTerm string

type LuceneEntry struct {
	ID     string
	Tokens []LuceneIndexTerm
}
