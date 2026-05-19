package query

import "github.com/LjungErik/zetra-lucene/lucene/search"

type Query interface {
	Execute(context search.SearchContext) []search.TopDoc
}
