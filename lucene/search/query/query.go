package query

import (
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
	"github.com/LjungErik/zetra-lucene/lucene/search/query/collector"
)

type Query interface {
	Execute(context.IndexReaderContext, collector.TopDocumentCollector)
}
