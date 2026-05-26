package index

import "iter"

type TermCount struct {
	DocumentID int
	Count      int
}

type Terms struct {
	terms []Term
}

func (t *Terms) Terms() iter.Seq[Term] {
	return func(yield func(Term) bool) {
		for _, term := range t.terms {
			if !yield(term) {
				return
			}
		}
	}
}

type Term struct {
	value    string
	postings []Posting
}

func (t *Term) Value() string {
	return t.value
}

func (t *Term) Postings() iter.Seq[Posting] {
	return func(yield func(Posting) bool) {
		for _, posting := range t.postings {
			if !yield(posting) {
				return
			}
		}
	}
}
