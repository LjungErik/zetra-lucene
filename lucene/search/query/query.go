package query

import (
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
)

type Query interface {
	Execute(ctx context.SearchContext) []document.TopDoc
}
